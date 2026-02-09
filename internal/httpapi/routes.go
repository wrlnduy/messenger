package httpapi

import (
	"messenger/internal/auth"
	"messenger/internal/cache"
	"messenger/internal/chat"
	"messenger/internal/sessions"
	"messenger/internal/storage"
	"messenger/internal/users"
	"messenger/internal/ws"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	Hub       *ws.Hub
	Store     storage.Store
	Auth      *auth.Service
	Users     users.Store
	Sessions  sessions.Store
	UserCache *cache.UserCache
}

func RegisterRoutes(mux *mux.Router, config *Config) {
	mux.Handle("/register", auth.RegisterHandler(config.Auth))

	mux.Handle("/login", auth.LoginHandler(config.Auth))

	mux.Handle("/logout", auth.LogoutHandler(config.Auth))

	logged := mux.PathPrefix("/logged").Subrouter()
	RegisterWithAuthRoutes(logged, config)

	mux.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))
}

func RegisterWithAuthRoutes(mux *mux.Router, config *Config) {
	mux.Use(func(next http.Handler) http.Handler {
		return auth.WithAuth(next, config.Auth)
	})

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(config.Hub, w, r)
	})

	mux.Handle("/message", chat.PostMessage(config.Hub, config.Store, config.UserCache))

	mux.HandleFunc("/history", chat.History(config.Store, config.UserCache))
}
