package auth

import (
	"context"
	"messenger/internal/users"

	"github.com/google/uuid"
)

type ctxKey string

const UserCtxKey ctxKey = "user"

func UserWithCtx(ctx context.Context) *users.User {
	user, _ := ctx.Value(UserCtxKey).(*users.User)
	return user
}

func UserIdWithCtx(ctx context.Context) uuid.UUID {
	user := UserWithCtx(ctx)
	userId, _ := uuid.Parse(*user.UserId)
	return userId
}
