package cache

import (
	"context"
	"fmt"
	"time"

	"messenger/internal/users"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	TTL = time.Hour
)

type UserCache struct {
	rdb   *redis.Client
	users users.Store
}

func NewUserCache(rdb *redis.Client, users users.Store) (*UserCache, error) {
	return &UserCache{rdb: rdb, users: users}, nil
}

func (c *UserCache) GetUsername(
	ctx context.Context,
	userID uuid.UUID,
) (string, error) {
	key := "user:" + userID.String()

	username, err := c.rdb.Get(ctx, key).Result()
	if err == nil {
		return username, nil
	}

	user, err := c.users.FindByID(ctx, userID)
	if err != nil {
		return "", err
	}

	username = *user.Username
	c.rdb.Set(ctx, key, username, TTL)

	return username, nil
}

func (c *UserCache) GetMapping(
	ctx context.Context,
	userIDs map[uuid.UUID]struct{},
) (map[uuid.UUID]string, error) {
	users := make(map[uuid.UUID]string)

	for userId := range userIDs {
		username, err := c.GetUsername(ctx, userId)
		if err != nil {
			return nil, err
		}

		users[userId] = username
	}

	return users, nil
}

func (c *UserCache) GetUserId(
	ctx context.Context,
	username string,
) (uuid.UUID, error) {
	key := fmt.Sprint("username:%v", username)

	id, err := c.rdb.Get(ctx, key).Result()
	if err == nil {
		uid, _ := uuid.Parse(id)
		return uid, nil
	}

	user, err := c.users.FindByUsername(ctx, username)
	if err != nil {
		return uuid.Nil, err
	}
	c.rdb.Set(ctx, key, *user.UserId, TTL)

	uid, _ := uuid.Parse(*user.UserId)
	return uid, nil
}
