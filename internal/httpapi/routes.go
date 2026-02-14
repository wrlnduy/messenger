package httpapi

import (
	"messenger/internal/auth"
	"messenger/internal/cache"
	"messenger/internal/chat"
	"messenger/internal/chats"
	"messenger/internal/messages"
	"messenger/internal/sessions"
	"messenger/internal/users"
	"messenger/internal/ws"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	Manager   *ws.HubManager
	Store     messages.Store
	Auth      *auth.Service
	Users     users.Store
	Chats     chats.Store
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

	mux.HandleFunc("/ws", ws.ServeWS(config.Manager))

	mux.Handle("/message", chat.PostMessage(config.Manager, config.Store, config.UserCache))

	mux.HandleFunc("/history", chat.History(config.Manager, config.Store, config.UserCache))

	mux.HandleFunc("/chats", chat.Chats(config.Chats, config.UserCache))

	mux.HandleFunc("/direct", chat.GetCreateDirect(config.Chats, config.UserCache))

	mux.HandleFunc("/group", chat.CreateGroup(config.Chats, config.UserCache))
}
