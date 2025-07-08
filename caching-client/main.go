package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shreyaskr1409/VidiMasta/caching-client/middlewares"
	pubsub "github.com/Shreyaskr1409/VidiMasta/caching-client/pub-sub"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "cache-client : ", log.LstdFlags)

	pubsub.Init(l)
	defer pubsub.Client.Close()
	l.Println("Valkey running at port 6379...")

	HTTPServerInit(l)
}

func HTTPServerInit(l *log.Logger) {
	mainRouter := mux.NewRouter()
	mainRouter.UseEncodedPath()
	mainRouter.Use(mux.CORSMethodMiddleware(mainRouter))
	mainRouter.Use(middlewares.LoggingMiddleware(l))

	s := &http.Server{
		Addr:         ":8082",
		Handler:      mainRouter,
		IdleTimeout:  120 * time.Second,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			l.Fatalln("Server encountered an error: ", err)
		}
	}()
	l.Println("HTTP Server listening at port 8082")

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGTERM)

	sig := <-shutdownChannel
	l.Println("Signal recieved for graceful shutdown: ", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	l.Println("Shutdown started...")
	s.Shutdown(ctx)
	l.Println("Shutdown successful!")
}
