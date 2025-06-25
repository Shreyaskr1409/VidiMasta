package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	l := log.New(os.Stdout, "ingestion-server: ", log.LstdFlags)
	l.Println("Logging starts")
}
