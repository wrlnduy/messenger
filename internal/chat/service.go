package chat

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"messenger/internal/chat/chats"
	"messenger/internal/chat/members"
	"messenger/internal/chat/messages"
	storemanager "messenger/internal/chat/store_manager"
	"messenger/internal/chat/stream"
	"messenger/internal/chat/unread"
	authpb "messenger/proto/auth"
	chatpb "messenger/proto/chats"
	userpb "messenger/proto/users"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	hubManager   *stream.HubManager
	storeManager storemanager.Manager
	messages     messages.Store
	chats        chats.Store
	members      members.Store
	unread       unread.Store

	usersClient userpb.UsersServiceClient
	authClient  authpb.AuthServiceClient
}

func NewService(
	hubManager *stream.HubManager,
	storeManager storemanager.Manager,
	messages messages.Store,
	chats chats.Store,
	members members.Store,
	unread unread.Store,
	usersClient userpb.UsersServiceClient,
	authClient authpb.AuthServiceClient,
) *Service {
	return &Service{
		hubManager:   hubManager,
		storeManager: storeManager,
		messages:     messages,
		chats:        chats,
		members:      members,
		unread:       unread,
		usersClient:  usersClient,
		authClient:   authClient,
	}
}

func (s *Service) GetHubByUser(
	ctx context.Context,
	userId uuid.UUID,
	chatId uuid.UUID,
) (*stream.Hub, error) {
	err := s.validateUser(ctx, userId, chatId)
	if err != nil {
		return nil, err
	}

	hub := s.hubManager.GetHub(chatId)
	return hub, nil
}

func (s *Service) PostMessage(
	ctx context.Context,
	hub *stream.Hub,
	chatId uuid.UUID,
	text string,
	userId uuid.UUID,
	username string,
) error {
	msg := &chatpb.ChatMessage{
		MessageId: proto.String(uuid.NewString()),
		UserId:    proto.String(userId.String()),
		ChatId:    proto.String(chatId.String()),
		Text:      proto.String(text),
		Timestamp: timestamppb.New(time.Now()),
	}

	err := s.storeManager.SaveMessage(ctx, msg)
	if err != nil {
		return err
	}

	msg.Username = proto.String(username)
	hub.Broadcast(userId, msg)

	return nil
}

func (s *Service) GetHistory(
	ctx context.Context,
	userId uuid.UUID,
	chatId uuid.UUID,
) (*chatpb.ChatHistory, error) {
	err := s.validateUser(ctx, userId, chatId)
	if err != nil {
		return nil, err
	}

	hist, err := s.messages.History(ctx, chatId)
	if err != nil {
		return nil, err
	}

	err = s.fillMapping(ctx, hist)
	if err != nil {
		log.Printf("Failed filling mapping for user:%q and chat:%q", userId, chatId)
		return nil, err
	}

	return hist, nil
}

func (s *Service) GetChats(
	ctx context.Context,
	userId uuid.UUID,
) (*chatpb.Chats, error) {
	chats, err := s.chats.GetUserChats(ctx, userId)
	if err != nil {
		return nil, err
	}

	for _, chat := range chats.Chats {
		if chat.Type.Number() != chatpb.ChatType_DIRECT.Number() {
			continue
		}

		chatId, _ := uuid.Parse(*chat.ChatId)
		buddyId, err := s.members.GetDirectBuddyId(ctx, chatId, userId)
		if err != nil {
			return nil, err
		}

		user, err := s.usersClient.FindByID(
			ctx,
			&userpb.FindByIDRequest{
				UserId: proto.String(buddyId.String()),
			},
		)
		if err != nil {
			return nil, err
		}
		chat.Title = user.Username
	}

	return chats, nil
}

func (s *Service) StartDirect(
	ctx context.Context,
	userId uuid.UUID,
	buddyName string,
) (*chatpb.Chat, error) {
	buddy, err := s.usersClient.FindByUsername(
		ctx,
		&userpb.FindByUsernameRequest{
			Username: proto.String(buddyName),
		},
	)
	if err != nil {
		return nil, err
	}
	buddyId := uuid.MustParse(buddy.GetUserId())

	chat, err := s.chats.GetDirect(ctx, userId, buddyId)
	if errors.Is(err, sql.ErrNoRows) {
		chat, err = s.storeManager.CreateDirect(ctx, userId, buddyId)
	}
	if err != nil {
		return nil, err
	}

	return chat, err
}

func (s *Service) StartGroup(
	ctx context.Context,
	creatorId uuid.UUID,
	title string,
	usernames []string,
) (*chatpb.Chat, error) {
	users, err := s.userIDsByUsernames(ctx, usernames)
	if err != nil {
		return nil, err
	}

	return s.storeManager.CreateGroup(ctx, creatorId, title, users)
}

func (s *Service) AddToGroup(
	ctx context.Context,
	userId uuid.UUID,
	chatId uuid.UUID,
) error {
	return s.storeManager.AddToGroup(ctx, userId, chatId)
}
