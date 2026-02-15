package stream

import (
	"messenger/internal/chat/chats"
	"sync"

	"github.com/google/uuid"
)

type HubManager struct {
	sync.Mutex
	hubs map[uuid.UUID]*Hub
}

func NewHubManager(chats chats.Store) *HubManager {
	return &HubManager{
		hubs: make(map[uuid.UUID]*Hub),
	}
}

func (m *HubManager) GetHub(chatId uuid.UUID) *Hub {
	m.Lock()
	defer m.Unlock()

	hub, ok := m.hubs[chatId]
	if !ok {
		hub = NewHub(chatId)
		m.hubs[chatId] = hub
	}

	return hub
}
