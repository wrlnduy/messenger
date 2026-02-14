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

func UserIDsByUsernames(ctx context.Context, cache *cache.UserCache, usernames []string) (uuid.UUIDs, error) {
	userIDs := make(map[uuid.UUID]struct{})
	for _, username := range usernames {
		id, err := cache.GetUserId(ctx, username)
		if err != nil {
			return nil, err
		}

		userIDs[id] = struct{}{}
	}

	ids := make(uuid.UUIDs, 0, len(userIDs))
	for id := range userIDs {
		ids = append(ids, id)
	}

	return ids, nil
}
