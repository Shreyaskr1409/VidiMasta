package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Shreyaskr1409/VidiMasta/gateway/middlewares"
	"github.com/Shreyaskr1409/VidiMasta/gateway/routes"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	l := log.New(os.Stdout, "auth: ", log.LstdFlags)
	l.Println("Logging starts")

	router := mux.NewRouter()
	router.UseEncodedPath()
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(middlewares.LoggingMiddleware(l))

	routes.HandleUserRoutes(router, l)

	conn, err := connectDB(l)
	if err != nil {
		l.Fatalln("Failed to connect to database: ", err)
	}
	defer conn.Close(context.Background())

	if err := conn.Ping(context.Background()); err != nil {
		l.Fatalln(err)
	}

	http.ListenAndServe(":9090", router)

	l.Println("Finished")
}

func connectDB(l *log.Logger) (*pgx.Conn, error) {
	err := godotenv.Load()
	if err != nil {
		l.Fatalln("Failed to obtain environment variables")
	}
	conn_str := os.Getenv("PGURL")
	return pgx.Connect(context.Background(), conn_str)
}
