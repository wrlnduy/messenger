package auth

import (
	"context"
	messenger "messenger/proto"

	"github.com/google/uuid"
)

type User *messenger.User

type ctxKey string

const UserCtxKey ctxKey = "user"

func UserWithCtx(ctx context.Context) User {
	user, _ := ctx.Value(UserCtxKey).(User)
	return user
}

func UserIdWithCtx(ctx context.Context) uuid.UUID {
	user := UserWithCtx(ctx)
	userId, _ := uuid.FromBytes([]byte(*user.UserId))
	return userId
}
