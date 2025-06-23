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

	userRouter.HandleFunc("/user", userHandler.GetUser).Methods("GET")
	userRouter.HandleFunc("/{username}", userHandler.GetUserByUsername).Methods("GET")
	userRouter.HandleFunc("/register", userHandler.Register).Methods("POST")
	userRouter.HandleFunc("/login", userHandler.Login).Methods("POST")
	userRouter.HandleFunc("/logout", userHandler.Logout).Methods("GET")
	userRouter.HandleFunc("/update", userHandler.UpdateUser).Methods("PATCH")
	userRouter.HandleFunc("/update-password", userHandler.UpdatePassword).Methods("PATCH")
}
