package handlers

import (
	"fmt"
	"log"
	"net/http"
)

type VideoHandler struct {
	l *log.Logger
}

func NewVideoHandler(l *log.Logger) *VideoHandler {
	vh := &VideoHandler{
		l: l,
	}
	return vh
}

func (h VideoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Working")
	fmt.Fprintln(w, "Hello World")

	if r.Method == http.MethodPost {
		h.uploadVideo(w, r)
	} else {
		http.Error(w, "Incorrect method, use POST", http.StatusBadRequest)
	}
}

func (h VideoHandler) uploadVideo(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Uploading Video")
	// file size limit for now is 10MB
	err := r.ParseMultipartForm(10<<20)
	if err != nil {
		http.Error(w, "Could not parse multipart form data", http.StatusBadRequest)
		return
	}
	h.l.Println("Parsed multipart form")
	// video-data is the name of the field
	file, handler, err := r.FormFile("video-data")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	h.l.Println("File name: ", handler.Filename)
	h.l.Println("File size: ", handler.Size)
	h.l.Println("File Headers: ", handler.Header)

	fmt.Fprintln(w, "Recieved File")
}
