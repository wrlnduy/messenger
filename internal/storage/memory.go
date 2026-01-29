package storage

import "sync"

type MemoryStore struct {
	sync.Mutex
	messages []Message
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) Save(m Message) {
	s.Lock()
	defer s.Unlock()

	s.messages = append(s.messages, m)
}

func (s *MemoryStore) List() []Message {
	s.Lock()
	defer s.Unlock()

	return append([]Message{}, s.messages...)
}
