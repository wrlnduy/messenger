package auth

import (
	"context"
	"errors"
	"log"
	"messenger/internal/auth/sessions"
	authpb "messenger/proto/auth"
	userpb "messenger/proto/users"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"
)

const (
	sessionTime = time.Hour * 24 * 30
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotActive      = errors.New("user not active")
)

type Service struct {
	sessions sessions.Store
	cache    *sessions.SessionsCache

	usersClient userpb.UsersServiceClient
}

func NewService(
	sessions sessions.Store,
	cache *sessions.SessionsCache,
	usersClient userpb.UsersServiceClient,
) *Service {
	return &Service{
		sessions:    sessions,
		cache:       cache,
		usersClient: usersClient,
	}
}

func (s *Service) Register(
	ctx context.Context,
	username string,
	password string,
) (*authpb.RegisterResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	resp, err := s.usersClient.CreateUser(
		ctx,
		&userpb.CreateUserRequest{
			Username:     proto.String(username),
			PasswordHash: proto.String(string(hash)),
		},
	)
	if err != nil {
		return nil, err
	}

	return &authpb.RegisterResponse{UserId: resp.UserId}, nil
}

func (s *Service) Login(
	ctx context.Context,
	username string,
	password string,
) (uuid.UUID, error) {
	user, err := s.usersClient.FindByUsername(
		ctx,
		&userpb.FindByUsernameRequest{
			Username: proto.String(username),
		},
	)
	if err != nil {
		return uuid.Nil, ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)) != nil {
		return uuid.Nil, ErrInvalidCredentials
	}

	if !*user.IsActive {
		return uuid.Nil, ErrUserNotActive
	}

	userId, _ := uuid.Parse(*user.UserId)
	sessionId, err := s.sessions.NewSession(ctx, userId, time.Now().Add(sessionTime))
	if err != nil {
		return uuid.Nil, err
	}

	return sessionId, nil
}

func (s *Service) Logout(
	ctx context.Context,
	sessionId uuid.UUID,
) error {
	return s.sessions.EndSession(ctx, sessionId)
}

func (s *Service) UserBySession(
	ctx context.Context,
	sessionId uuid.UUID,
) (*userpb.User, error) {
	log.Printf("Searching for user by sessionId: %v\n", sessionId)

	userId, err := s.cache.UserByID(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	return s.usersClient.FindByID(
		ctx,
		&userpb.FindByIDRequest{
			UserId: proto.String(userId.String()),
		},
	)
}
