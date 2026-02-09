package chat

import (
	"context"
	"messenger/internal/cache"
	messenger "messenger/proto"

	"github.com/google/uuid"
)

func FillMapping(ctx context.Context, hist *messenger.ChatHistory, cache *cache.UserCache) error {
	userIDs := make(map[uuid.UUID]struct{})
	for _, msg := range hist.Messages {
		id, _ := uuid.Parse(*msg.UserId)
		userIDs[id] = struct{}{}
	}

	users, err := cache.GetMapping(ctx, userIDs)
	if err != nil {
		return err
	}

	hist.Mapping = make(map[string]string)
	for userId, username := range users {
		hist.Mapping[userId.String()] = username
	}

	return nil
}
