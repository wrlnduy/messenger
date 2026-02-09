package db

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRdb(param string) (*redis.Client, error) {
	opt, err := redis.ParseURL(os.Getenv(param))
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
