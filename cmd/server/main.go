package main

import (
	"log"
	"net/http"
	"os"

	"messenger/internal/auth"
	"messenger/internal/cache"
	"messenger/internal/db"
	"messenger/internal/httpapi"
	"messenger/internal/messages"
	"messenger/internal/sessions"
	"messenger/internal/users"
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

	dbs, err := db.NewDb(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbs.Close()

	rdb, err := db.NewRdb("REDIS_URL")
	if err != nil {
		log.Fatal(err)
	}
	defer rdb.Close()

	users, err := users.NewPostgresStore(dbs)
	if err != nil {
		log.Fatal(err)
	}

	store, err := messages.NewPostgresStore(dbs)
	if err != nil {
		log.Fatal(err)
	}

	sessions, err := sessions.NewPostgresStore(dbs)
	if err != nil {
		log.Fatal(err)
	}

	auth, err := auth.NewService(dbs, users, sessions)
	if err != nil {
		log.Fatal(err)
	}

	userCache, err := cache.NewUserCache(rdb, users)
	if err != nil {
		log.Fatal(err)
	}

	mux := mux.NewRouter()
	httpapi.RegisterRoutes(mux, &httpapi.Config{
		Hub:       hub,
		Store:     store,
		Auth:      auth,
		Users:     users,
		Sessions:  sessions,
		UserCache: userCache,
	})

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
