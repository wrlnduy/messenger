package users

import (
	"context"
	"fmt"
	userpb "messenger/proto/users"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

const (
	TTL = time.Hour
)

type UserCache struct {
	rdb   *redis.Client
	users Store
}

func NewUserCache(rdb *redis.Client, users Store) (*UserCache, error) {
	return &UserCache{rdb: rdb, users: users}, nil
}

func (c *UserCache) FindByID(
	ctx context.Context,
	userId uuid.UUID,
) (*userpb.User, error) {
	key := "user:" + userId.String()

	user := &userpb.User{}

	bytes, err := c.rdb.Get(ctx, key).Bytes()
	if err == nil {
		err = proto.Unmarshal(bytes, user)
		return user, nil
	}

	user, err = c.users.FindByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	data, _ := proto.Marshal(user)
	c.rdb.Set(ctx, key, data, TTL)

	return user, nil
}

func (c *UserCache) GetMapping(
	ctx context.Context,
	userIDs uuid.UUIDs,
) (map[uuid.UUID]string, error) {
	users := make(map[uuid.UUID]string)

	for _, userId := range userIDs {
		user, err := c.FindByID(ctx, userId)
		if err != nil {
			return nil, err
		}
		users[userId] = *user.Username
	}

	return users, nil
}

func (c *UserCache) FindByUsername(
	ctx context.Context,
	username string,
) (*userpb.User, error) {
	key := fmt.Sprint("username:%v", username)

	user := &userpb.User{}
	bytes, err := c.rdb.Get(ctx, key).Bytes()
	if err == nil {
		proto.Unmarshal(bytes, user)
		return user, nil
	}

	user, err = c.users.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	data, _ := proto.Marshal(user)
	c.rdb.Set(ctx, key, data, TTL)

	return user, nil
}
