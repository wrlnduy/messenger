package users

import (
	"context"
	userpb "messenger/proto/users"

	"github.com/google/uuid"
)

type Service struct {
	users      Store
	usersCache *UserCache
}

func NewService(
	users Store,
	usersCache *UserCache,
) *Service {
	return &Service{
		users:      users,
		usersCache: usersCache,
	}
}

func (s *Service) CreateUser(
	ctx context.Context,
	username, passwordHash string,
) (uuid.UUID, error) {
	userId := uuid.New()
	err := s.users.CreateUser(ctx, userId, username, passwordHash)
	return userId, err
}

func (s *Service) FindByID(
	ctx context.Context,
	userId uuid.UUID,
) (*userpb.User, error) {
	return s.usersCache.FindByID(ctx, userId)
}

func (s *Service) FindByUsername(
	ctx context.Context,
	username string,
) (*userpb.User, error) {
	return s.usersCache.FindByUsername(ctx, username)
}

func (s *Service) GetMapping(
	ctx context.Context,
	userIDs uuid.UUIDs,
) (map[uuid.UUID]string, error) {
	return s.usersCache.GetMapping(ctx, userIDs)
}
