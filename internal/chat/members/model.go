package members

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Store interface {
	IsMember(ctx context.Context, chatId uuid.UUID, userId uuid.UUID) (bool, error)

	GetDirectBuddyId(ctx context.Context, chatId uuid.UUID, userId uuid.UUID) (uuid.UUID, error)

	AddMembers(ctx context.Context, tx *sql.Tx, chatId uuid.UUID, userIDs uuid.UUIDs) error

	UpdateRole(ctx context.Context, tx *sql.Tx, chatId uuid.UUID, userId uuid.UUID, role string) error
}
