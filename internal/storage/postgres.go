package storage

import (
	"context"
	"database/sql"
	message "messenger/proto"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/protobuf/proto"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) (*PostgresStore, error) {
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Save(ctx context.Context, msg *message.ChatMessage) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO chat_messages(message_id, user_id, text, timestamp) VALUES($1,$2,$3,$4)`,
		msg.MessageId, msg.UserId, msg.Text, time.Unix(*msg.Timestamp, 0),
	)
	return err
}

func (s *PostgresStore) List(ctx context.Context) ([]*message.ChatMessage, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT message_id, user_id, text, timestamp FROM chat_messages ORDER BY timestamp ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*message.ChatMessage
	for rows.Next() {
		var msg message.ChatMessage
		var t time.Time
		if err := rows.Scan(&msg.MessageId, &msg.UserId, &msg.Text, &t); err != nil {
			return nil, err
		}
		msg.Timestamp = proto.Int64(t.Unix())
		messages = append(messages, &msg)
	}
	return messages, rows.Err()
}

func (s *PostgresStore) History(ctx context.Context) (*message.ChatHistory, error) {
	msgs, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	return &message.ChatHistory{Messages: msgs}, nil
}
