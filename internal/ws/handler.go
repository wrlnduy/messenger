package ws

import (
	"messenger/internal/cookies"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 30 * time.Second,
	CheckOrigin:      func(r *http.Request) bool { return true },
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := NewClient(hub, conn, cookies.UserID(r))
	hub.register <- client

	go client.ReadPump()
	go client.WritePump()
}
