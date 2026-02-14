package ws

import (
	"github.com/google/uuid"
)

type Hub struct {
	ChatId uuid.UUID

	clients    map[uuid.UUID]*Client
	register   chan *Client
	unregister chan *Client
}

func NewHub(chatId uuid.UUID) *Hub {
	return &Hub{
		ChatId:     chatId,
		clients:    make(map[uuid.UUID]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c.UserID] = c
		case c := <-h.unregister:
			if _, ok := h.clients[c.UserID]; ok {
				delete(h.clients, c.UserID)
				close(c.send)
			}
		}
	}
}

func (h *Hub) Broadcast(msg []byte) {
	for _, c := range h.clients {
		c.PostMessage(msg)
	}
}
