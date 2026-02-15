package sessions

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	TTL = time.Hour
)

type SessionsCache struct {
	rdb      *redis.Client
	sessions Store
}

func NewSessionsCache(rdb *redis.Client, sessions Store) *SessionsCache {
	return &SessionsCache{
		rdb:      rdb,
		sessions: sessions,
	}
}

func (c *SessionsCache) UserByID(
	ctx context.Context,
	sessionId uuid.UUID,
) (uuid.UUID, error) {
	key := sessionKey(sessionId)

	id, err := c.rdb.Get(ctx, key).Result()
	if err == nil {
		return uuid.MustParse(id), nil
	}

	userId, err := c.sessions.UserByID(ctx, sessionId)
	if err != nil {
		return uuid.Nil, err
	}
	fmt.Printf("Setting {%q, %q}\n", key, userId.String())
	c.rdb.Set(ctx, key, userId.String(), TTL)

	return userId, nil
}

func sessionKey(sessionId uuid.UUID) string {
	return fmt.Sprintf("session:%v", sessionId)
}
