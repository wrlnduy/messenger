package storage

import (
	"time"
)

type Message struct {
	UserID    string
	Text      string
	CreatedAt time.Time
}

type Store interface {
	Save(Message)
	List() []Message
}
