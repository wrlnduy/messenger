package messages

import (
	"context"
	messenger "messenger/proto"

	"github.com/google/uuid"
)

type Store interface {
	Save(context.Context, *messenger.ChatMessage) error
	List(context.Context, uuid.UUID) ([]*messenger.ChatMessage, error)
	History(context.Context, uuid.UUID) (*messenger.ChatHistory, error)
}
