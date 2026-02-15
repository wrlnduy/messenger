package messages

import (
	"context"
	"database/sql"
	chatpb "messenger/proto/chats"

	"github.com/google/uuid"
)

type Store interface {
	Save(ctx context.Context, tx *sql.Tx, msg *chatpb.ChatMessage) error

	List(ctx context.Context, userId uuid.UUID) ([]*chatpb.ChatMessage, error)

	History(ctx context.Context, userId uuid.UUID) (*chatpb.ChatHistory, error)
}
