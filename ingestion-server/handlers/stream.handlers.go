package handlers

import (
	"bytes"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Shreyaskr1409/VidiMasta/ingestion-server/data"
	"github.com/nareix/joy4/format/rtmp"
)

type StreamHandler struct {
	l           *log.Logger
	storagePath string        // Should end with '/'
	mtx         *sync.RWMutex // Ensures that a transcoding happens only once
	streams     map[string]*data.Stream
}

func NewStreamHandler(l *log.Logger, storagePath string, mtx *sync.RWMutex) *StreamHandler {
	return &StreamHandler{
		l:           l,
		storagePath: storagePath,
		mtx:         mtx,
		streams:     make(map[string]*data.Stream),
	}
}

func (h *StreamHandler) Publish(conn *rtmp.Conn) {
	streamKey := strings.TrimPrefix(conn.URL.Path, "/")
	println(streamKey)

	cmd := exec.Command("ffmpeg",
		// Input options
		"-f", "flv", // Force FLV input format
		"-analyzeduration", "10M", // Increase probe size
		"-probesize", "10M", // Increase analysis duration
		"-i", "pipe:0", // Read from stdin

		// Audio handling (if needed)
		"-acodec", "aac", // Transcode audio to AAC
		"-ar", "44100", // Set audio sample rate
		"-b:a", "128k", // Set audio bitrate

		// Video handling
		"-vcodec", "copy", // Copy video stream as-is

		// HLS output options
		"-f", "hls",
		"-hls_time", "2",
		"-hls_list_size", "10",
		"-hls_flags", "delete_segments",

		// Proper path joining
		filepath.Join(h.storagePath, streamKey+".m3u8"),
	)
	cmd.Stdin = conn.NetConn()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	h.mtx.Lock()
	h.streams[streamKey] = &data.Stream{Cmd: cmd}
	h.mtx.Unlock()

	err := cmd.Run() // Block until FFmpeg exits
	if err != nil {
		h.l.Printf("FFmpeg error: %v\n", err)
		h.l.Printf("FFmpeg stdout: %s\n", stdout.String())
		h.l.Printf("FFmpeg stderr: %s\n", stderr.String())
	}

	defer func() {
		h.mtx.Lock()
		delete(h.streams, streamKey)
		h.mtx.Unlock()
	}() // Ensures cleanup even if FFmpeg crashes
}
