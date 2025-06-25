package main

import (
	"log"
	"os"
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
	sh := handlers.NewStreamHandler(l, "/public/storage/", &sync.RWMutex{})

	server := &rtmp.Server{
		HandlePublish: sh.Publish,
		HandlePlay:    nil,
		HandleConn:    nil,
	}

	l.Println("Server is listening at port :8081")
	err := server.ListenAndServe()
	if err != nil {
		l.Fatal(err)
	}
}
