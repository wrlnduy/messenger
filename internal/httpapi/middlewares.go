package httpapi

import (
	"context"
	"messenger/internal/contextx"
	"messenger/internal/cookies"
	authpb "messenger/proto/auth"
	"net/http"

	"google.golang.org/protobuf/proto"
)

func (g *Gateway) withAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId, ok := cookies.SessionID(r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := g.authClient.UserBySession(
			r.Context(),
			&authpb.UserBySessionRequest{
				SessionId: proto.String(sessionId.String()),
			},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), contextx.UserCtxKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
