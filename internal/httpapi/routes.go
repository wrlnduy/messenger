package httpapi

import (
	"encoding/json"
	"messenger/internal/chat"
	"messenger/internal/storage"
	"messenger/internal/ws"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, hub *ws.Hub, store storage.Store) {
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(hub, w, r)
	})

	mux.Handle("/message", chat.PostMessage(hub, store))

	mux.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		msgs := store.List()
		json.NewEncoder(w).Encode(msgs)
	})

	mux.Handle("/", http.FileServer(http.Dir("./web")))
}
