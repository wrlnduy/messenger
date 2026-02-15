package users

import (
	"context"
	"database/sql"
	"errors"

	userpb "messenger/proto/users"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

var (
	ErrorNotAdmin = errors.New("User should be admin")
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) (*PostgresStore, error) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS users (
			user_id UUID PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,

			is_active BOOLEAN NOT NULL DEFAULT TRUE,
			is_admin BOOLEAN NOT NULL DEFAULT FALSE,

			created_at TIMESTAMP NOT NULL DEFAULT now()
		);`,
	)
	if err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) CreateUser(ctx context.Context, userId uuid.UUID, username, passwordHash string) error {
	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO users (user_id, username, password_hash) VALUES ($1, $2, $3)`,
		userId, username, passwordHash,
	)
	return err
}

func (s *PostgresStore) FindByUsername(ctx context.Context, username string) (*userpb.User, error) {
	user := &userpb.User{
		Username: proto.String(username),
	}

	err := s.db.QueryRowContext(
		ctx,
		`SELECT user_id, password_hash, is_active, is_admin FROM users WHERE username = $1`,
		username,
	).Scan(&user.UserId, &user.PasswordHash, &user.IsActive, &user.IsAdmin)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *PostgresStore) FindByID(ctx context.Context, id uuid.UUID) (*userpb.User, error) {
	user := &userpb.User{
		UserId: proto.String(id.String()),
	}

	err := s.db.QueryRowContext(
		ctx,
		`SELECT username, password_hash, is_active, is_admin FROM users WHERE user_id = $1`,
		id,
	).Scan(&user.Username, &user.PasswordHash, &user.IsActive, &user.IsAdmin)

	if err != nil {
		return nil, err
	}
	return user, nil
}
