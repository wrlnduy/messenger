package main

import (
	"log"
	"net/http"
	"os"

	"messenger/internal/auth"
	"messenger/internal/db"
	"messenger/internal/httpapi"
	"messenger/internal/storage"
	"messenger/internal/ws"

	"github.com/gorilla/mux"
)

func main() {
	hub := ws.NewHub()
	go hub.Run()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := db.New(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store, err := storage.NewPostgresStore(db)
	if err != nil {
		log.Fatal(err)
	}

	auth, err := auth.NewService(db)
	if err != nil {
		log.Fatal(err)
	}

	mux := mux.NewRouter()
	httpapi.RegisterRoutes(mux, &httpapi.Config{
		Hub:   hub,
		Store: store,
		Auth:  auth,
	})

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
