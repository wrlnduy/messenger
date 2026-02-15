package usergrpc

import (
	"context"
	"messenger/internal/users"
	userpb "messenger/proto/users"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	userpb.UnimplementedUsersServiceServer
	service *users.Service
}

func NewServer(service *users.Service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) CreateUser(
	ctx context.Context,
	req *userpb.CreateUserRequest,
) (*userpb.CreateUserResponse, error) {
	userId, err := s.service.CreateUser(ctx, req.GetUsername(), req.GetPasswordHash())
	if err != nil {
		return nil, err
	}

	return &userpb.CreateUserResponse{
		UserId: proto.String(userId.String()),
	}, nil
}

func (s *Server) FindByID(
	ctx context.Context,
	req *userpb.FindByIDRequest,
) (*userpb.User, error) {
	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, err
	}

	return s.service.FindByID(ctx, userId)
}

func (s *Server) FindByUsername(
	ctx context.Context,
	req *userpb.FindByUsernameRequest,
) (*userpb.User, error) {
	return s.service.FindByUsername(ctx, req.GetUsername())
}

func (s *Server) GetMapping(
	ctx context.Context,
	req *userpb.GetMappingRequest,
) (*userpb.GetMappingResponse, error) {
	usedIDs := make(uuid.UUIDs, 0, len(req.GetUserIDs()))
	for _, id := range req.GetUserIDs() {
		userId, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}

		usedIDs = append(usedIDs, userId)
	}

	mapping, err := s.service.GetMapping(ctx, usedIDs)
	if err != nil {
		return nil, err
	}

	smapping := make(map[string]string, len(usedIDs))
	for k, v := range mapping {
		smapping[k.String()] = v
	}

	return &userpb.GetMappingResponse{
		Mapping: smapping,
	}, nil
}
