package main

import (
	"log"
	"net/http"

	"messenger/internal/httpapi"
	"messenger/internal/storage"
	"messenger/internal/ws"
)

func main() {
	hub := ws.NewHub()
	go hub.Run()

	store := storage.NewMemoryStore()

	mux := http.NewServeMux()
	httpapi.RegisterRoutes(mux, hub, store)

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
