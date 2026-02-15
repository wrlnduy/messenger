package sessions

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) (*PostgresStore, error) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS sessions (
			session_id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(user_id),
			
			expires_at TIMESTAMP NOT NULL
		);`,
	)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) NewSession(
	ctx context.Context,
	userId uuid.UUID,
	expiresAt time.Time,
) (uuid.UUID, error) {
	sessionId := uuid.New()

	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO sessions (session_id, user_id, expires_at) VALUES ($1, $2, $3)`,
		sessionId, userId, expiresAt,
	)

	return sessionId, err
}

func (s *PostgresStore) EndSession(
	ctx context.Context,
	sessionId uuid.UUID,
) error {
	_, err := s.db.ExecContext(
		ctx,
		`UPDATE sessions SET expires_at = $1 WHERE session_id = $2`,
		time.Now(), sessionId,
	)
	return err
}

func (s *PostgresStore) UserByID(
	ctx context.Context,
	sessionId uuid.UUID,
) (uuid.UUID, error) {
	var userId uuid.UUID
	err := s.db.QueryRowContext(
		ctx, `
		SELECT user_id
		FROM sessions
		WHERE session_id = $1
			AND expires_at > now()
		`,
		sessionId,
	).Scan(&userId)

	return userId, err
}
