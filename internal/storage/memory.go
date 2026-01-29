package storage

import (
	"messenger/proto"
	"sync"
)

type MemoryStore struct {
	sync.Mutex
	messages []*messenger.ChatMessage
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) Save(m *messenger.ChatMessage) {
	s.Lock()
	defer s.Unlock()

	s.messages = append(s.messages, m)
}

func (s *MemoryStore) List() []*messenger.ChatMessage {
	s.Lock()
	defer s.Unlock()

	return append([]*messenger.ChatMessage{}, s.messages...)
}

func (s *MemoryStore) History() *messenger.ChatHistory {
	return &messenger.ChatHistory{
		Messages: s.List(),
	}
}
