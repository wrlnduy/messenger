package auth

import (
	"context"
)

type ctxKey string

const UserCtxKey ctxKey = "user"

func UserWithCtx(ctx context.Context) {
	// user, _ := ctx.Value(UserCtxKey)
}
