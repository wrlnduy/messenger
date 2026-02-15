package storemanager

import (
	"context"
	"database/sql"
	"errors"
	"messenger/internal/chat/chats"
	"messenger/internal/chat/members"
	"messenger/internal/chat/messages"
	"messenger/internal/chat/unread"
	chatpb "messenger/proto/chats"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostgresStoreManager struct {
	db       *sql.DB
	messages messages.Store
	chats    chats.Store
	members  members.Store
	unread   unread.Store
}

func NewPostgresStore(
	db *sql.DB,
	messages messages.Store,
	chats chats.Store,
	members members.Store,
	unread unread.Store,
) Manager {
	return &PostgresStoreManager{
		db,
		messages,
		chats,
		members,
		unread,
	}
}

func (m *PostgresStoreManager) SaveMessage(
	ctx context.Context,
	msg *chatpb.ChatMessage,
) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = m.messages.Save(ctx, tx, msg)
	if err != nil {
		return err
	}

	err = m.unread.IncrementUnread(
		ctx,
		tx,
		uuid.MustParse(*msg.ChatId),
		uuid.MustParse(*msg.UserId),
	)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *PostgresStoreManager) CreateDirect(
	ctx context.Context,
	u1, u2 uuid.UUID,
) (*chatpb.Chat, error) {
	if u1 == u2 {
		return nil, errors.New("cannot create direct chat with yourself")
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	chatId := uuid.New()
	chat := &chatpb.Chat{
		ChatId:    proto.String(chatId.String()),
		Type:      chatpb.ChatType_DIRECT.Enum(),
		CreatedAt: timestamppb.Now(),
	}

	m.chats.CreateDirect(ctx, tx, chat)

	m.members.AddMembers(ctx, tx, chatId, []uuid.UUID{u1, u2})

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return chat, nil
}

func (m *PostgresStoreManager) CreateGroup(
	ctx context.Context,
	creator uuid.UUID,
	title string,
	users uuid.UUIDs,
) (*chatpb.Chat, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	chatId := uuid.New()
	chat := &chatpb.Chat{
		ChatId:    proto.String(chatId.String()),
		Type:      chatpb.ChatType_GROUP.Enum(),
		Title:     proto.String(title),
		CreatedAt: timestamppb.Now(),
	}

	err = m.chats.CreateGroup(ctx, tx, chat)
	if err != nil {
		return nil, err
	}

	users = append(users, creator)
	err = m.members.AddMembers(ctx, tx, chatId, users)
	if err != nil {
		return nil, err
	}

	err = m.members.UpdateRole(ctx, tx, chatId, creator, chatpb.ChatRole_ADMIN.String())

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return chat, nil
}

func (m *PostgresStoreManager) AddToGroup(
	ctx context.Context,
	userId uuid.UUID,
	chatId uuid.UUID,
) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = m.members.AddMembers(ctx, tx, chatId, []uuid.UUID{userId})
	if err != nil {
		return err
	}

	return tx.Commit()
}
