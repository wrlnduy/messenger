package auth

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"messenger/internal/users"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionTime = time.Hour * 24 * 30
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotActive      = errors.New("user not active")
)

type Service struct {
	db    *sql.DB
	users users.Store
}

func NewService(db *sql.DB, users users.Store) (*Service, error) {
	return &Service{
		db:    db,
		users: users,
	}, nil
}

func (s *Service) Register(
	ctx context.Context,
	username string,
	password string,
) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.users.CreateUser(ctx, uuid.New(), username, string(hash))
}

func (s *Service) Login(
	ctx context.Context,
	username string,
	password string,
) (uuid.UUID, error) {
	user, err := s.users.FindByUsername(ctx, username)

	if err != nil {
		return uuid.Nil, ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)) != nil {
		return uuid.Nil, ErrInvalidCredentials
	}

	if !*user.IsActive {
		return uuid.Nil, ErrUserNotActive
	}

	sessionId := uuid.New()

	_, err = s.db.ExecContext(
		ctx,
		`INSERT INTO sessions (session_id, user_id, expires_at) VALUES ($1, $2, $3)`,
		sessionId, *user.UserId, time.Now().Add(sessionTime),
	)

	return sessionId, err
}

func (s *Service) UserBySession(
	ctx context.Context,
	sessionId uuid.UUID,
) (*users.User, error) {
	user := new(users.User)

	log.Printf("Searching for user by sessionId: %v\n", sessionId)
	err := s.db.QueryRowContext(
		ctx,
		`SELECT u.user_id, u.username, u.is_active, u.is_admin
		FROM sessions s
		JOIN users u ON u.user_id = s.user_id
		WHERE s.session_id = $1`,
		sessionId,
	).Scan(&user.UserId, &user.Username, &user.IsActive, &user.IsAdmin)

	if err != nil {
		return nil, err
	}

	return user, err
}
