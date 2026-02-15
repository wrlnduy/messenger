package authgrpc

import (
	"context"
	"messenger/internal/auth"
	authpb "messenger/proto/auth"
	userpb "messenger/proto/users"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	authpb.UnimplementedAuthServiceServer
	service *auth.Service
}

func NewServer(service *auth.Service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) Register(
	ctx context.Context,
	req *authpb.RegisterRequest,
) (*authpb.RegisterResponse, error) {
	return s.service.Register(
		ctx,
		req.GetUsername(),
		req.GetPassword(),
	)
}

func (s *Server) Login(
	ctx context.Context,
	req *authpb.LoginRequest,
) (*authpb.LoginResponse, error) {
	sessionId, err := s.service.Login(
		ctx,
		req.GetUsername(),
		req.GetPassword(),
	)
	if err != nil {
		return nil, err
	}

	return &authpb.LoginResponse{
		SessionId: proto.String(sessionId.String()),
	}, nil
}

func (s *Server) Logout(
	ctx context.Context,
	req *authpb.LogoutRequest,
) (*emptypb.Empty, error) {
	sessionId, err := uuid.Parse(req.GetSessionId())
	if err != nil {
		return nil, err
	}

	err = s.service.Logout(ctx, sessionId)
	return &emptypb.Empty{}, err
}

func (s *Server) UserBySession(
	ctx context.Context,
	req *authpb.UserBySessionRequest,
) (*userpb.User, error) {
	sessionId, err := uuid.Parse(req.GetSessionId())
	if err != nil {
		return nil, err
	}

	return s.service.UserBySession(ctx, sessionId)
}
