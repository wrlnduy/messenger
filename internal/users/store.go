package users

import (
	"context"
	"database/sql"
	"errors"

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
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) CreateUser(ctx context.Context, id uuid.UUID, username, passwordHash string) error {
	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO users (user_id, username, password_hash) VALUES ($1, $2, $3)`,
		id, username, passwordHash,
	)
	return err
}

func (s *PostgresStore) FindByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{
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

func (s *PostgresStore) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	user := &User{
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

func (s *PostgresStore) ApproveUser(ctx context.Context, approval_id uuid.UUID, id uuid.UUID) error {
	admin, err := s.FindByID(ctx, approval_id)
	if err != nil {
		return err
	}

	if !*admin.IsAdmin {
		return ErrorNotAdmin
	}

	_, err = s.FindByID(ctx, id)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `UPDATE users SET is_active = true WHERE user_id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}
