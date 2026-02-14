package ws

import (
	"messenger/internal/auth"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 30 * time.Second,
	CheckOrigin:      func(r *http.Request) bool { return true },
}

func ServeWS(manager *HubManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hub, err := manager.GetHubForRequst(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		client := NewClient(hub, conn, auth.UserIdWithCtx(r.Context()))
		hub.register <- client

		go client.ReadPump()
		go client.WritePump()
	}
}
