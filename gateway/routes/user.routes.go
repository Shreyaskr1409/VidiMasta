package routes

import (
	"log"

	"github.com/Shreyaskr1409/VidiMasta/gateway/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleUserRoutes(router *mux.Router, l *log.Logger, db *pgxpool.Pool) {
	userRouter := router.PathPrefix("/api/v1/users").Subrouter()
	userHandler := handlers.NewUserHandler(l, db)

	userRouter.HandleFunc("/register", userHandler.Register).Methods("POST")
}
