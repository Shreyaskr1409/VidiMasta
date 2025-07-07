package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/Shreyaskr1409/VidiMasta/ingestion-server/database"
	"github.com/Shreyaskr1409/VidiMasta/ingestion-server/handlers"
	"github.com/Shreyaskr1409/VidiMasta/ingestion-server/middlewares"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nareix/joy4/format"
	"github.com/nareix/joy4/format/rtmp"
)

func main() {
	godotenv.Load()
	l := log.New(os.Stdout, "ingestion-server: ", log.LstdFlags)
	l.Println("Logging starts")

	handleRTMP(l)

	database.Init(l)
	defer database.Close()
	err := database.Migrate(database.DB, l, context.Background())
	if err != nil {
		l.Fatalln("Errpr encountered during database migration: ", err)
	}

	router := mux.NewRouter()
	router.UseEncodedPath()
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(middlewares.LoggingMiddleware(l))

	s := &http.Server{
		Addr:         ":8081",
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			l.Fatal(err)
		}
	}()
	l.Println("Server is listening at port :8080")

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGTERM)

	sig := <-shutdownChannel
	l.Println("Recieved signal for graceful shutdown. signal: ", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	l.Println("Shutdown started")
	s.Shutdown(ctx)
	l.Println("Shutdown successful")
}

func handleRTMP(l *log.Logger) {
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
