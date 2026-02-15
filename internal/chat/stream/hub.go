package stream

import (
	chatpb "messenger/proto/chats"
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	sync.Mutex

	ChatId  uuid.UUID
	clients map[uuid.UUID]chatpb.ChatService_ConnectServer
}

func NewHub(chatId uuid.UUID) *Hub {
	return &Hub{
		ChatId:  chatId,
		clients: make(map[uuid.UUID]chatpb.ChatService_ConnectServer),
	}
}

func (h *Hub) Register(userId uuid.UUID, stream chatpb.ChatService_ConnectServer) {
	h.Lock()
	defer h.Unlock()
	h.clients[userId] = stream
}

func (h *Hub) Unregister(userId uuid.UUID) {
	h.Lock()
	defer h.Unlock()
	delete(h.clients, userId)
}

func (h *Hub) Broadcast(sender uuid.UUID, msg *chatpb.ChatMessage) {
	h.Lock()
	defer h.Unlock()

	for _, conn := range h.clients {
		conn.Send(msg)
	}
}
