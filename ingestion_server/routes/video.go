package routes

import (
	"log"
	"net/http"

	"github.com/Shreyaskr1409/VidiMasta/ingestion_server/handlers"
)

func HandleVideoRoutes(l *log.Logger) *http.ServeMux {
	videoMux := http.NewServeMux()
	videoHandler := handlers.NewVideoHandler(l)
	videoMux.Handle("/", videoHandler)
	return videoMux
}
