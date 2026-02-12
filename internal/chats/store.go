package chats

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	messenger "messenger/proto"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) (*PostgresStore, error) {
	_, err := db.Exec(
		`DO $$
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

		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'chat_role') THEN
				CREATE TYPE chat_role AS ENUM ('MEMBER', 'ADMIN');
			END IF;
		END
		$$;

		CREATE TABLE IF NOT EXISTS chat_members (
			chat_id UUID NOT NULL REFERENCES chats(chat_id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
			role chat_role DEFAULT 'member',
			PRIMARY KEY (chat_id, user_id)
		);`,
	)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) GetByID(
	ctx context.Context,
	chatId uuid.UUID,
) (*Chat, error) {
	chat := new(Chat)

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
	chat.Type = messenger.ChatType(messenger.ChatType_value[ctype]).Enum()

	return chat, nil
}

func (s *PostgresStore) GetUserChats(
	ctx context.Context,
	userId uuid.UUID,
) (*Chats, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT chat_id, user_id FROM chat_members WHERE user_id = $1`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chats := make([]*Chat, 0)

	for rows.Next() {
		var chatId uuid.UUID
		err = rows.Scan(&chatId, &userId)
		if err != nil {
			return nil, err
		}

		chat, err := s.GetByID(ctx, chatId)
		if err != nil {
			return nil, err
		}

		chats = append(chats, chat)
	}

	return &Chats{Chats: chats}, nil
}

func (s *PostgresStore) GetDirect(
	ctx context.Context,
	u1, u2 uuid.UUID,
) (*Chat, error) {
	var chatId uuid.UUID
	err := s.db.QueryRowContext(ctx,
		`SELECT c.chat_id
		FROM chats c
		JOIN chat_members cm1 ON cm1.chat_id = c.chat_id
		JOIN chat_members cm2 ON cm2.chat_id = c.chat_id
		WHERE c.type = 'DIRECT'
		AND cm1.user_id = $1
		AND cm2.user_id = $2`,
		u1, u2,
	).Scan(&chatId)
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, chatId)
}

func (s *PostgresStore) CreateDirect(
	ctx context.Context,
	u1, u2 uuid.UUID,
) (*Chat, error) {
	if u1 == u2 {
		return nil, errors.New("cannot create direct chat with yourself")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	chat := &Chat{
		ChatId:    proto.String(uuid.NewString()),
		Type:      messenger.ChatType_DIRECT.Enum(),
		CreatedAt: timestamppb.Now(),
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO chats (chat_id, type, created_at) VALUES ($1, $2, $3)`,
		*chat.ChatId,
		chat.Type.String(),
		chat.CreatedAt.AsTime(),
	)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2), ($1, $3)`,
		*chat.ChatId,
		u1,
		u2,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *PostgresStore) CreateGroup(
	ctx context.Context,
	creator uuid.UUID,
	title string,
	users uuid.UUIDs,
) (*Chat, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	chat := &Chat{
		ChatId:    proto.String(uuid.NewString()),
		Type:      messenger.ChatType_GROUP.Enum(),
		Title:     proto.String(title),
		CreatedAt: timestamppb.Now(),
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO chats (chat_id, type, created_at) VALUES ($1, $2, $3)`,
		*chat.ChatId,
		chat.Type.String(),
		chat.CreatedAt.AsTime(),
	)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO chat_members (chat_id, user_id, role) VALUES ($1, $2, $3)`,
		*chat.ChatId,
		creator,
		"ADMIN",
	)
	if err != nil {
		return nil, err
	}

	if len(users) > 0 {
		var (
			q    bytes.Buffer
			args = make([]any, 0, len(users)+1)
		)
		q.WriteString(`INSERT INTO chat_members (chat_id, user_id) VALUES `)
		args = append(args, *chat.ChatId)

		for i, u := range users {
			if i > 0 {
				q.WriteString(`, `)
			}

			q.WriteString(fmt.Sprint("($1, $%v)", i+2))
			args = append(args, u)
		}

		_, err = tx.ExecContext(ctx, q.String(), args...)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *PostgresStore) IsMember(
	ctx context.Context,
	chatId uuid.UUID,
	userId uuid.UUID,
) (bool, error) {
	var exists int
	err := s.db.QueryRowContext(ctx,
		`SELECT 1
		FROM chat_members
		WHERE chat_id = $1 AND user_id = $2`,
		chatId,
		userId,
	).Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *PostgresStore) GetDirectBuddyId(
	ctx context.Context,
	chatId uuid.UUID,
	userId uuid.UUID,
) (uuid.UUID, error) {
	var uId uuid.UUID
	err := s.db.QueryRowContext(ctx,
		`SELECT c.user_id
		FROM chat_members c
		WHERE chat_id = $1 AND user_id != $2`,
		chatId,
		userId,
	).Scan(&uId)
	return uId, err
}
