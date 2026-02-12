package chats

import (
	"context"
	messenger "messenger/proto"

	"github.com/google/uuid"
)

var (
	GlobalChatID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
)

type Chat = messenger.Chat

type Chats = messenger.Chats

type Store interface {
	GetByID(ctx context.Context, chatId uuid.UUID) (*Chat, error)

	GetUserChats(ctx context.Context, userId uuid.UUID) (*Chats, error)

	GetDirect(ctx context.Context, u1, u2 uuid.UUID) (*Chat, error)

	CreateDirect(ctx context.Context, u1, u2 uuid.UUID) (*Chat, error)

	CreateGroup(ctx context.Context, creator uuid.UUID, title string, users uuid.UUIDs) (*Chat, error)

	IsMember(ctx context.Context, chatId uuid.UUID, userId uuid.UUID) (bool, error)
}
