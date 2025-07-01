package handlers

import (
	"bytes"
	"log"
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

	cmdConfig := data.TranscodingCommand{
		Input: data.InputOptions{
			Format:          "flv",
			AnalyseDuration: "10M",
			ProbeSize:       "10M",
			Input:           "pipe:0",
		},
		Audio: data.AudioOptions{
			Codec:      "aac",
			SampleRate: "44100",
			AudioRate:  "128k",
		},
		Video: data.VideoOptions{
			Codec: "copy",
		},
		Output: data.OutputOptions{
			Format:   "hls",
			Time:     "2",
			ListSize: "10",
			Flags:    "delete_segments",
		},
	}

	// Generate the command
	cmd := data.BuildCommand(cmdConfig, h.storagePath, streamKey)
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
