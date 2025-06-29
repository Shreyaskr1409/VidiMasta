package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/Shreyaskr1409/VidiMasta/ingestion-server/handlers"
	"github.com/joho/godotenv"
	"github.com/nareix/joy4/format"
	"github.com/nareix/joy4/format/rtmp"
)

func main() {
	godotenv.Load()
	l := log.New(os.Stdout, "ingestion-server: ", log.LstdFlags)
	l.Println("Logging starts")

	format.RegisterAll()
	sh := handlers.NewStreamHandler(l, getAbsPath("./public/storage"), &sync.RWMutex{})

	server := &rtmp.Server{
		HandlePublish: sh.Publish,
		HandlePlay:    nil,
		HandleConn:    nil,
	}

	err := server.ListenAndServe()
	if err != nil {
		l.Fatal(err)
	}
}

func getAbsPath(s string) string {
	relativePath := s

	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		log.Fatal("Failed to get absolute path:", err)
	}

	// Ensure trailing slash (optional, but useful for directories)
	absPath = filepath.Join(absPath, "") // Adds `/` at the end

	log.Println("Absolute Path:", absPath)

	return absPath
}
