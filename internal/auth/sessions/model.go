package sessions

import (
	"context"
	"time"

	authpb "messenger/proto/auth"

	"github.com/google/uuid"
)

type Session authpb.Session

type Store interface {
	NewSession(ctx context.Context, userId uuid.UUID, expiresAt time.Time) (uuid.UUID, error)

	EndSession(ctx context.Context, sessionId uuid.UUID) error

	UserByID(ctx context.Context, sessionId uuid.UUID) (uuid.UUID, error)
}
