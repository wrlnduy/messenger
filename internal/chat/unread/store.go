package unread

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(
	db *sql.DB,
) (Store, error) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS chat_unread (
			chat_id UUID NOT NULL REFERENCES chats(chat_id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
			unread_count INT NOT NULL DEFAULT 0,
			PRIMARY KEY (chat_id, user_id)
		);
	`)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) IncrementUnread(
	ctx context.Context,
	tx *sql.Tx,
	chatId uuid.UUID,
	writer uuid.UUID,
) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO chat_unread (chat_id, user_id, unread_count)
		SELECT cm.chat_id, cm.user_id, 1
		FROM chat_members cm
		WHERE cm.chat_id = $1
		  AND cm.user_id != $2
		ON CONFLICT (chat_id, user_id)
		DO UPDATE SET unread_count = chat_unread.unread_count + 1
	`,
		chatId,
		writer,
	)

	return err
}

func (s *PostgresStore) ResetUnread(
	ctx context.Context,
	chatId uuid.UUID,
	userId uuid.UUID,
) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE chat_unread
		SET unread_count = 0
		WHERE chat_id = $1 AND user_id = $2
	`,
		chatId,
		userId,
	)

	return err
}
