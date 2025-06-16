package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleUserRoutes(router *mux.Router, l *log.Logger) {
	userRouter := router.PathPrefix("/api/v1/users").Subrouter()
	userRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	})
}
