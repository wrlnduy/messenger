package storage

import (
	"context"
	message "messenger/proto"
)

type Store interface {
	Save(context.Context, *message.ChatMessage) error
	List(context.Context) ([]*message.ChatMessage, error)
	History(context.Context) (*message.ChatHistory, error)
}
