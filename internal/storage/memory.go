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

func (s *MemoryStore) Save(m *message.ChatMessage) error {
	s.Lock()
	defer s.Unlock()

	s.messages = append(s.messages, m)
	return nil
}

func (s *MemoryStore) List() ([]*message.ChatMessage, error) {
	s.Lock()
	defer s.Unlock()

	return append([]*message.ChatMessage{}, s.messages...), nil
}

func (s *MemoryStore) History() (*message.ChatHistory, error) {
	msgs, _ := s.List()
	return &message.ChatHistory{
		Messages: msgs,
	}, nil
}
