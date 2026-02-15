package users

import (
	"context"
	userpb "messenger/proto/users"

	"github.com/google/uuid"
)

type Store interface {
	CreateUser(ctx context.Context, userId uuid.UUID, username, passwordHash string) error

	FindByUsername(ctx context.Context, username string) (*userpb.User, error)

	FindByID(ctx context.Context, id uuid.UUID) (*userpb.User, error)
}
