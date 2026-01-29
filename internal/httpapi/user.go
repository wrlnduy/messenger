package httpapi

import (
	"net/http"
	"time"

	"messenger/internal/cookies"

	"github.com/google/uuid"
)

func WithUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie(cookies.UserIDCookie)
		if err == http.ErrNoCookie {
			id := uuid.NewString()

			http.SetCookie(w, &http.Cookie{
				Name:     cookies.UserIDCookie,
				Value:    id,
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Now().Add(365 * 24 * time.Hour),
			})

			r.AddCookie(&http.Cookie{
				Name:  cookies.UserIDCookie,
				Value: id,
			})
		}

		next.ServeHTTP(w, r)
	})
}
