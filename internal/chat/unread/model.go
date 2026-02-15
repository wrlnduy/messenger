package unread

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Store interface {
	IncrementUnread(ctx context.Context, tx *sql.Tx, chatId uuid.UUID, writer uuid.UUID) error

	ResetUnread(ctx context.Context, chatId uuid.UUID, userId uuid.UUID) error
}
