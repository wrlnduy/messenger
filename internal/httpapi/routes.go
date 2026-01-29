package httpapi

import (
	"messenger/internal/chat"
	"messenger/internal/storage"
	"messenger/internal/ws"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
)

func RegisterRoutes(mux *http.ServeMux, hub *ws.Hub, store storage.Store) {
	handler := http.NewServeMux()

	handler.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(hub, w, r)
	})

	handler.Handle("/message", chat.PostMessage(hub, store))

	handler.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		data, _ := protojson.Marshal(store.History())
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	handler.Handle("/", http.FileServer(http.Dir("./web")))

	mux.Handle("/", WithUser(handler))
}
