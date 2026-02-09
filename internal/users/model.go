package users

import (
	"context"
	messenger "messenger/proto"

	"github.com/google/uuid"
)

type User messenger.User

type Store interface {
	CreateUser(ctx context.Context, id uuid.UUID, username, passwordHash string) error

	FindByUsername(ctx context.Context, username string) (*User, error)

	FindByID(ctx context.Context, id uuid.UUID) (*User, error)

	ApproveUser(ctx context.Context, approval_id uuid.UUID, id uuid.UUID) error
}
