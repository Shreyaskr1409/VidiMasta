package handlers

import (
	"log"
	"os/exec"
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
	streamKey := conn.URL.Path

	cmd := exec.Command("ffmpeg",
		"-i", "pipe:0", // Reading from stdin
		"-c", "copy", // No re-encoding
		"-f", "hls", // sets format
		"-hls_time", "2", // 2-second segments
		"-hls_list_size", "10", // Keep 10 segments (20s DVR)
		"-hls_flags", "delete_segments", // Delete old segments
		h.storagePath+streamKey+".m3u8", // Output playlist
	)
	cmd.Stdin = conn.NetConn() // Sending input from connection to the standard input

	h.mtx.Lock()
	h.streams[streamKey].Cmd = cmd
	h.mtx.Unlock()

	err := cmd.Run() // Block until FFmpeg exits
	if err != nil {
		h.l.Println("Error encountered while transcoding: \n", err)
	}

	defer func() {
		h.mtx.Lock()
		delete(h.streams, streamKey)
		h.mtx.Unlock()
	}() // Ensures cleanup even if FFmpeg crashes
}
