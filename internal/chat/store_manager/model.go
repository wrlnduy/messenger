package storemanager

import (
	"context"
	chatpb "messenger/proto/chats"

	"github.com/google/uuid"
)

type Manager interface {
	SaveMessage(ctx context.Context, msg *chatpb.ChatMessage) error

	CreateDirect(ctx context.Context, u1, u2 uuid.UUID) (*chatpb.Chat, error)

	CreateGroup(ctx context.Context, creator uuid.UUID, title string, users uuid.UUIDs) (*chatpb.Chat, error)

	AddToGroup(ctx context.Context, userId uuid.UUID, chatId uuid.UUID) error
}
