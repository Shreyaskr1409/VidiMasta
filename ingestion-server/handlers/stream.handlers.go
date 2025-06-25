package handlers

import (
	"log"
	"net"
	"net/http"
)

type StreamHandler struct {
	l           *log.Logger
	storagePath string
}

func NewStreamHandler(l *log.Logger, storagePath string) *StreamHandler {
	return &StreamHandler{
		storagePath: storagePath,
	}
}

func (h *StreamHandler) ServeHTTP(w http.ResponseWriter, r *http.Response)

func (h *StreamHandler) HandleRTMPConn(conn net.Conn)

func (h *StreamHandler) uploadToStorage(localPath string, quality string)
