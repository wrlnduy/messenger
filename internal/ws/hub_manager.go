package ws

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/auth"
	"messenger/internal/chats"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type HubManager struct {
	sync.Mutex
	hubs  map[uuid.UUID]*Hub
	chats chats.Store
}

func NewHubManager(chats chats.Store) *HubManager {
	return &HubManager{
		hubs:  make(map[uuid.UUID]*Hub),
		chats: chats,
	}
}

func (m *HubManager) GetHub(chatId uuid.UUID) *Hub {
	m.Lock()
	defer m.Unlock()

	hub, ok := m.hubs[chatId]
	if !ok {
		hub = NewHub(chatId)
		m.hubs[chatId] = hub
		go hub.Run()
	}

	return hub
}

func (m *HubManager) GetHubByUser(
	ctx context.Context,
	userId uuid.UUID,
	chatId uuid.UUID,
) (*Hub, error) {
	ok, err := m.chats.IsMember(ctx, chatId, userId)
	if err != nil || !ok {
		return nil, errors.New(fmt.Sprintf("User %q should member of chat %q", userId, chatId))
	}

	hub := m.GetHub(chatId)
	return hub, nil
}

func (m *HubManager) GetHubForRequst(
	r *http.Request,
) (*Hub, error) {
	chatIdStr := r.URL.Query().Get("chat_id")
	chatId, err := uuid.Parse(chatIdStr)
	if err != nil {
		return nil, err
	}

	userId := auth.UserIdWithCtx(r.Context())

	return m.GetHubByUser(r.Context(), userId, chatId)
}
