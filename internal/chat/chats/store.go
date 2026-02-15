package chats

import (
	"context"
	"database/sql"
	chatpb "messenger/proto/chats"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) (*PostgresStore, error) {
	_, err := db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'chat_type') THEN
				CREATE TYPE chat_type AS ENUM ('GLOBAL', 'DIRECT', 'GROUP');
			END IF;
		END
		$$;

		CREATE TABLE IF NOT EXISTS chats (
			chat_id UUID PRIMARY KEY,
			type chat_type NOT NULL,
			title TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT now()
		);
		`,
	)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) GetByID(
	ctx context.Context,
	chatId uuid.UUID,
) (*chatpb.Chat, error) {
	chat := &chatpb.Chat{}

	var t time.Time
	var ctype string
	err := s.db.QueryRowContext(ctx,
		`SELECT chat_id, type, title, created_at
		FROM chats WHERE chat_id = $1`,
		chatId,
	).Scan(&chat.ChatId, &ctype, &chat.Title, &t)
	if err != nil {
		return nil, err
	}
	chat.CreatedAt = timestamppb.New(t)
	chat.Type = chatpb.ChatType(chatpb.ChatType_value[ctype]).Enum()

	return chat, nil
}

func (s *PostgresStore) CreateDirect(
	ctx context.Context,
	tx *sql.Tx,
	chat *chatpb.Chat,
) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO chats (chat_id, type, created_at) VALUES ($1, $2, $3)`,
		*chat.ChatId,
		chat.Type.String(),
		chat.CreatedAt.AsTime(),
	)
	return err
}

func (s *PostgresStore) CreateGroup(
	ctx context.Context,
	tx *sql.Tx,
	chat *chatpb.Chat,
) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO chats (chat_id, type, created_at, title) VALUES ($1, $2, $3, $4)`,
		*chat.ChatId,
		chat.Type.String(),
		chat.CreatedAt.AsTime(),
		*chat.Title,
	)
	return err
}

func (s *PostgresStore) GetUserChats(
	ctx context.Context,
	userId uuid.UUID,
) (*chatpb.Chats, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT 
			c.chat_id,
			c.title,
			c.type,
			COALESCE(cu.unread_count, 0) as unread
		FROM chats c
		JOIN chat_members cm ON cm.chat_id = c.chat_id
		LEFT JOIN chat_unread cu 
			ON cu.chat_id = c.chat_id 
		AND cu.user_id = $1
		WHERE cm.user_id = $1;
		`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chats := make([]*chatpb.Chat, 0)

	for rows.Next() {
		var chat chatpb.Chat
		var ctype string
		err = rows.Scan(&chat.ChatId, &chat.Title, &ctype, &chat.Unread)
		if err != nil {
			return nil, err
		}
		chat.Type = chatpb.ChatType(chatpb.ChatType_value[ctype]).Enum()

		chats = append(chats, &chat)
	}

	return &chatpb.Chats{Chats: chats}, nil
}

func (s *PostgresStore) GetDirect(
	ctx context.Context,
	u1, u2 uuid.UUID,
) (*chatpb.Chat, error) {
	var chatId uuid.UUID
	err := s.db.QueryRowContext(ctx, `
		SELECT c.chat_id
		FROM chats c
		JOIN chat_members cm1 ON cm1.chat_id = c.chat_id
		JOIN chat_members cm2 ON cm2.chat_id = c.chat_id
		WHERE c.type = 'DIRECT'
		AND cm1.user_id = $1
		AND cm2.user_id = $2
		`,
		u1, u2,
	).Scan(&chatId)
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, chatId)
}
