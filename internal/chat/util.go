package chat

import (
	"context"
	"errors"
	"fmt"

	chatpb "messenger/proto/chats"
	userpb "messenger/proto/users"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func (s *Service) fillMapping(ctx context.Context, hist *chatpb.ChatHistory) error {
	uniqueUserIDs := make(map[uuid.UUID]struct{})
	for _, msg := range hist.Messages {
		uniqueUserIDs[uuid.MustParse(*msg.UserId)] = struct{}{}
	}

	userIDs := make(uuid.UUIDs, 0, len(uniqueUserIDs))
	for id := range uniqueUserIDs {
		userIDs = append(userIDs, id)
	}

	resp, err := s.usersClient.GetMapping(ctx, &userpb.GetMappingRequest{
		UserIDs: userIDs.Strings(),
	})
	if err != nil {
		return err
	}

	hist.Mapping = resp.Mapping

	return nil
}

func (s *Service) userIDsByUsernames(ctx context.Context, usernames []string) (uuid.UUIDs, error) {
	userIDs := make(map[uuid.UUID]struct{})
	for _, username := range usernames {
		resp, err := s.usersClient.FindByUsername(
			ctx,
			&userpb.FindByUsernameRequest{
				Username: proto.String(username),
			},
		)
		if err != nil {
			return nil, err
		}

		userIDs[uuid.MustParse(resp.GetUserId())] = struct{}{}
	}

	ids := make(uuid.UUIDs, 0, len(userIDs))
	for id := range userIDs {
		ids = append(ids, id)
	}

	return ids, nil
}

func (s *Service) validateUser(
	ctx context.Context,
	userId uuid.UUID,
	chatId uuid.UUID,
) error {
	ok, err := s.members.IsMember(ctx, chatId, userId)
	if err != nil || !ok {
		return errors.New(fmt.Sprintf("User %q should member of chat %q", userId, chatId))
	}

	return nil
}

func markMine(
	msgs []*chatpb.ChatMessage,
	userId uuid.UUID,
) {
	id := userId.String()
	for _, msg := range msgs {
		msg.IsMine = proto.Bool(msg.GetUserId() == id)
	}
}
