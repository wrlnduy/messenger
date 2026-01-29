package main

import (
	"log"
	"net/http"
	"os"

	"messenger/internal/httpapi"
	"messenger/internal/storage"
	"messenger/internal/ws"
)

func main() {
	hub := ws.NewHub()
	go hub.Run()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	store, err := storage.NewPostgresStore(dsn)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	httpapi.RegisterRoutes(mux, hub, store)

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
