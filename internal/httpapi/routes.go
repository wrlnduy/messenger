package httpapi

import (
	"messenger/internal/auth"
	"messenger/internal/chat"
	"messenger/internal/storage"
	"messenger/internal/ws"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/protobuf/encoding/protojson"
)

type Config struct {
	Hub   *ws.Hub
	Store storage.Store
	Auth  *auth.Service
}

func RegisterRoutes(mux *mux.Router, config *Config) {
	mux.Handle("/register", auth.RegisterHandler(config.Auth))

	mux.Handle("/login", auth.LoginHandler(config.Auth))

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

	mux.Handle("/message", chat.PostMessage(config.Hub, config.Store))

	mux.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		hist, err := config.Store.History(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, _ := protojson.Marshal(hist)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
}
