package contextx

import (
	"context"
	userpb "messenger/proto/users"

	"github.com/google/uuid"
)

type ctxKey string

const UserCtxKey ctxKey = "user"

func UserWithCtx(ctx context.Context) *userpb.User {
	user, _ := ctx.Value(UserCtxKey).(*userpb.User)
	return user
}

func UserIdWithCtx(ctx context.Context) uuid.UUID {
	user := UserWithCtx(ctx)
	return uuid.MustParse(*user.UserId)
}
