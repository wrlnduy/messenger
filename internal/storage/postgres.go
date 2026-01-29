package storage

import (
	"context"
	"database/sql"

	"messenger/proto"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(dsn string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS chat_messages (
		message_id UUID PRIMARY KEY,
		user_id UUID NOT NULL,
		text TEXT NOT NULL,
		timestamp BIGINT NOT NULL
	)
	`)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Save(ctx context.Context, msg *message.ChatMessage) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO chat_messages(message_id, user_id, text, timestamp) VALUES($1,$2,$3,$4)`,
		msg.MessageId, msg.UserId, msg.Text, msg.Timestamp,
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
		if err := rows.Scan(&msg.MessageId, &msg.UserId, &msg.Text, &msg.Timestamp); err != nil {
			return nil, err
		}
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
