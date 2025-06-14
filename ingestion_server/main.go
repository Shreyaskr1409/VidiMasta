package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Shreyaskr1409/VidiMasta/ingestion_server/routes"
)

func main() {
	l := log.New(os.Stdout, "ingestion server: ", log.LstdFlags)
	l.Println("Hello World")

	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World");
	})
	mainMux.Handle("/video", routes.HandleVideoRoutes(l))

	err := http.ListenAndServe(":4646", mainMux)
	if err != nil {
		l.Println("Error encountered while starting the server")
	}
}
