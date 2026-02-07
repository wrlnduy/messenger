package auth

import (
	"context"
	"database/sql"
	"errors"
	"log"
	messenger "messenger/proto"
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
	db *sql.DB
}

func NewService(db *sql.DB) (*Service, error) {
	return &Service{
		db: db,
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

	_, err = s.db.ExecContext(ctx,
		`INSERT INTO users (user_id, username, password_hash) VALUES ($1, $2, $3)`,
		uuid.New(),
		username,
		string(hash),
	)
	return err
}

func (s *Service) Login(
	ctx context.Context,
	username string,
	password string,
) (uuid.UUID, error) {
	var (
		userId   uuid.UUID
		hash     string
		isActive bool
	)

	err := s.db.QueryRowContext(
		ctx,
		`SELECT user_id, password_hash, is_active FROM users WHERE username = $1`,
		username,
	).Scan(&userId, &hash, &isActive)

	if err != nil {
		return uuid.Nil, ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return uuid.Nil, ErrInvalidCredentials
	}

	if !isActive {
		return uuid.Nil, ErrUserNotActive
	}

	sessionId := uuid.New()

	_, err = s.db.ExecContext(
		ctx,
		`INSERT INTO sessions (session_id, user_id, expires_at) VALUES ($1, $2, $3)`,
		sessionId, userId, time.Now().Add(sessionTime),
	)

	return sessionId, err
}

func (s *Service) UserBySession(
	ctx context.Context,
	sessionId uuid.UUID,
) (User, error) {
	var u messenger.User

	log.Printf("Searching for user by sessionId: %v\n", sessionId)
	err := s.db.QueryRowContext(
		ctx,
		`SELECT u.user_id, u.username, u.is_active, u.is_admin
		FROM sessions s
		JOIN users u ON u.user_id = s.user_id
		WHERE s.session_id = $1`,
		sessionId,
	).Scan(&u.UserId, &u.Username, &u.IsActive, &u.IsAdmin)

	if err != nil {
		return nil, err
	}

	return &u, err
}
