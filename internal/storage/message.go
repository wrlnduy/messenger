package storage

import (
	"context"
	"messenger/proto"
)

type Store interface {
	Save(*message.ChatMessage) error
	List() ([]*message.ChatMessage, error)
	History() (*message.ChatHistory, error)
}

type StoreContext interface {
	Save(context.Context, *message.ChatMessage) error
	List(context.Context) ([]*message.ChatMessage, error)
	History(context.Context) (*message.ChatHistory, error)
}
