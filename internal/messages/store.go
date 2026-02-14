package messages

import (
	"context"
	"database/sql"
	messenger "messenger/proto"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) (*PostgresStore, error) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS chat_messages (
			message_id UUID PRIMARY KEY,
			chat_id UUID NOT NULL REFERENCES chats(chat_id) ON DELETE CASCADE,
			
			user_id UUID NOT NULL REFERENCES users(user_id),
			text TEXT NOT NULL,
			timestamp TIMESTAMP NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_chat_messages_chat_id_timestamp
		ON chat_messages (chat_id, timestamp);`,
	)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Save(ctx context.Context, msg *messenger.ChatMessage) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO chat_messages(message_id, chat_id, user_id, text, timestamp) VALUES($1, $2, $3, $4, $5)`,
		*msg.MessageId, *msg.ChatId, *msg.UserId, *msg.Text, msg.Timestamp.AsTime(),
	)
	return err
}

func (s *PostgresStore) List(ctx context.Context, chat_id uuid.UUID) ([]*messenger.ChatMessage, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT message_id, user_id, text, timestamp 
		FROM chat_messages c
		WHERE c.chat_id = $1
		ORDER BY timestamp ASC`,
		chat_id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*messenger.ChatMessage, 0)
	for rows.Next() {
		var msg messenger.ChatMessage
		var t time.Time
		if err := rows.Scan(&msg.MessageId, &msg.UserId, &msg.Text, &t); err != nil {
			return nil, err
		}
		msg.Timestamp = timestamppb.New(t)
		messages = append(messages, &msg)
	}
	return messages, rows.Err()
}

func (s *PostgresStore) History(ctx context.Context, chat_id uuid.UUID) (*messenger.ChatHistory, error) {
	msgs, err := s.List(ctx, chat_id)
	if err != nil {
		return nil, err
	}
	return &messenger.ChatHistory{Messages: msgs}, nil
}
