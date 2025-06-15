package main

import (
	"log"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "auth: ", log.LstdFlags)
	l.Println("Logging starts")

	router := mux.NewRouter()
	router.UseEncodedPath()

	userRouter := router.PathPrefix("/api/v1/users").Subrouter()
}
