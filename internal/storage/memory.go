package storage

import (
	"messenger/proto"
	"sync"
)

type MemoryStore struct {
	sync.Mutex
	messages []*message.ChatMessage
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) Save(m *message.ChatMessage) {
	s.Lock()
	defer s.Unlock()

	s.messages = append(s.messages, m)
}

func (s *MemoryStore) List() []*message.ChatMessage {
	s.Lock()
	defer s.Unlock()

	return append([]*message.ChatMessage{}, s.messages...)
}

func (s *MemoryStore) History() *message.ChatHistory {
	return &message.ChatHistory{
		Messages: s.List(),
	}
}
