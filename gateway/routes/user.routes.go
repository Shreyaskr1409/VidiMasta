package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleUserRoutes(router *mux.Router, l *log.Logger, db *pgxpool.Pool) {
	userRouter := router.PathPrefix("/api/v1/users").Subrouter()
	userRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	})
}
