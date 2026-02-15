package chats

import (
	"context"
	"database/sql"
	chatpb "messenger/proto/chats"

	"github.com/google/uuid"
)

var (
	GlobalChatID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
)

type Store interface {
	GetByID(ctx context.Context, chatId uuid.UUID) (*chatpb.Chat, error)

	CreateDirect(ctx context.Context, tx *sql.Tx, chat *chatpb.Chat) error

	CreateGroup(ctx context.Context, tx *sql.Tx, chat *chatpb.Chat) error

	GetUserChats(ctx context.Context, userId uuid.UUID) (*chatpb.Chats, error)

	GetDirect(ctx context.Context, u1, u2 uuid.UUID) (*chatpb.Chat, error)
}
