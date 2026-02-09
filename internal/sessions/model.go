package sessions

import (
	"context"
	messenger "messenger/proto"
	"time"

	"github.com/google/uuid"
)

type Session messenger.Session

type Store interface {
	NewSession(ctx context.Context, userId uuid.UUID, expiresAt time.Time) (uuid.UUID, error)

	EndSession(ctx context.Context, sessionId uuid.UUID) error
}
