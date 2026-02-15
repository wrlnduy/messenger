package chatgrpc

import (
	"context"
	"log"
	"messenger/internal/chat"
	"messenger/internal/chat/chats"
	chatpb "messenger/proto/chats"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	chatpb.UnimplementedChatServiceServer
	service *chat.Service
}

func NewServer(service *chat.Service) *Server {
	return &Server{service: service}
}

func (s *Server) Connect(
	stream chatpb.ChatService_ConnectServer,
) error {
	firstMsg, err := stream.Recv()
	if err != nil {
		return err
	}

	chatId := uuid.MustParse(firstMsg.GetChatId())
	userId := uuid.MustParse(firstMsg.GetUserId())
	username := firstMsg.GetUsername()

	hub, err := s.service.GetHubByUser(stream.Context(), userId, chatId)
	if err != nil {
		return err
	}

	hub.Register(userId, stream)
	defer hub.Unregister(userId)

	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}

		err = s.service.PostMessage(
			stream.Context(),
			hub,
			chatId,
			msg.GetText(),
			userId,
			username,
		)
		if err != nil {
			return err
		}
	}
}

func (s *Server) GetHistory(
	ctx context.Context,
	req *chatpb.HistoryRequest,
) (*chatpb.ChatHistory, error) {
	chatId, err := uuid.Parse(req.GetChatId())
	if err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, err
	}

	log.Printf("Getting history for (%q, %q)\n", chatId, userId)

	return s.service.GetHistory(ctx, userId, chatId)
}

func (s *Server) GetChats(
	ctx context.Context,
	req *chatpb.GetChatsRequest,
) (*chatpb.Chats, error) {
	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, err
	}
	return s.service.GetChats(ctx, userId)
}

func (s *Server) StartDirect(
	ctx context.Context,
	req *chatpb.StartDirectRequest,
) (*chatpb.Chat, error) {
	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, err
	}
	return s.service.StartDirect(ctx, userId, req.GetBuddy())
}

func (s *Server) StartGroup(
	ctx context.Context,
	req *chatpb.StartGroupRequest,
) (*chatpb.Chat, error) {
	creatorId, err := uuid.Parse(req.GetCreatorId())
	if err != nil {
		return nil, err
	}
	return s.service.StartGroup(ctx, creatorId, req.GetTitle(), req.GetUsernames())
}

func (s *Server) AddToGroup(
	ctx context.Context,
	req *chatpb.AddToGroupRequest,
) (*emptypb.Empty, error) {
	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, err
	}

	if req.GetChatId() == "" {
		return &emptypb.Empty{}, s.service.AddToGroup(ctx, userId, chats.GlobalChatID)
	}

	chatId, err := uuid.Parse(req.GetChatId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, s.service.AddToGroup(ctx, userId, chatId)
}
