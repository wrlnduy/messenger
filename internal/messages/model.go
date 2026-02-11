package messages

import (
	"context"
	messenger "messenger/proto"
)

type Store interface {
	Save(context.Context, *messenger.ChatMessage) error
	List(context.Context) ([]*messenger.ChatMessage, error)
	History(context.Context) (*messenger.ChatHistory, error)
}
