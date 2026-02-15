package members

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) (Store, error) {
	_, err := db.Exec(`
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
		);
	`)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) IsMember(
	ctx context.Context,
	chatId uuid.UUID,
	userId uuid.UUID,
) (bool, error) {
	var exists int
	err := s.db.QueryRowContext(ctx, `
		SELECT 1
		FROM chat_members
		WHERE chat_id = $1 AND user_id = $2
		`,
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
	err := s.db.QueryRowContext(ctx, `
		SELECT c.user_id
		FROM chat_members c
		WHERE chat_id = $1 AND user_id != $2
		`,
		chatId,
		userId,
	).Scan(&uId)
	return uId, err
}

func (s *PostgresStore) AddMembers(
	ctx context.Context,
	tx *sql.Tx,
	chatId uuid.UUID,
	userIDs uuid.UUIDs,
) error {
	var (
		q    bytes.Buffer
		args = make([]any, 0, len(userIDs)+1)
	)
	q.WriteString(`INSERT INTO chat_members (chat_id, user_id) VALUES `)
	args = append(args, chatId)

	for i, u := range userIDs {
		if i > 0 {
			q.WriteString(`, `)
		}

		q.WriteString(fmt.Sprintf("($1, $%v)", i+2))
		args = append(args, u)
	}

	_, err := tx.ExecContext(ctx, q.String(), args...)
	return err
}

func (s *PostgresStore) UpdateRole(
	ctx context.Context,
	tx *sql.Tx,
	chatId uuid.UUID,
	userId uuid.UUID,
	role string,
) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE chat_members
		SET role = $1
		WHERE chat_id = $2
			AND user_id = $3 
	`,
		role,
		chatId,
		userId,
	)
	return err
}
