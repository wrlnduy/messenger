package messages

import (
	"context"
	"database/sql"
	messenger "messenger/proto"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/protobuf/proto"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) (*PostgresStore, error) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS chat_messages (
    message_id UUID PRIMARY KEY,
    chat_id UUID NOT NULL REFERENCES chats(chat_id) ON DELETE CASCADE,
    
    user_id UUID NOT NULL REFERENCES users(user_id),
    text TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL
	);`)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Save(ctx context.Context, msg *messenger.ChatMessage) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO chat_messages(message_id, user_id, text, timestamp) VALUES($1,$2,$3,$4)`,
		msg.MessageId, msg.UserId, msg.Text, time.Unix(*msg.Timestamp, 0),
	)
	return err
}

func (s *PostgresStore) List(ctx context.Context) ([]*messenger.ChatMessage, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT message_id, user_id, text, timestamp FROM chat_messages ORDER BY timestamp ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*messenger.ChatMessage
	for rows.Next() {
		var msg messenger.ChatMessage
		var t time.Time
		if err := rows.Scan(&msg.MessageId, &msg.UserId, &msg.Text, &t); err != nil {
			return nil, err
		}
		msg.Timestamp = proto.Int64(t.Unix())
		messages = append(messages, &msg)
	}
	return messages, rows.Err()
}

func (s *PostgresStore) History(ctx context.Context) (*messenger.ChatHistory, error) {
	msgs, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	return &messenger.ChatHistory{Messages: msgs}, nil
}
