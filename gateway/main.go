package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	l := log.New(os.Stdout, "auth: ", log.LstdFlags)
	l.Println("Logging starts")

	router := mux.NewRouter()
	router.UseEncodedPath()
	router.Use(mux.CORSMethodMiddleware(router))

	// userRouter := router.PathPrefix("/api/v1/users").Subrouter()
	db, err := connectDB(l)
	if err != nil {
		l.Fatalln("Failed to connect to database: ", err)
	}

	if err := db.Ping(); err != nil {
		l.Fatalln(err)
	}

	l.Println("Finished")
}

func connectDB(l *log.Logger) (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		l.Fatalln("Failed to load env")
	}
	db_username := os.Getenv("PGUSER")
	db_password := os.Getenv("PGPASSWORD")

	conn_str := fmt.Sprintf("postgresql://%s:%s@ep-still-hat-a9l6pp6h-pooler.gwc.azure.neon.tech/neondb?sslmode=require", db_username, db_password)
	l.Println(conn_str)

	return sql.Open("postgres", conn_str)
}
