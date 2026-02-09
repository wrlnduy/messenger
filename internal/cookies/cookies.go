package cookies

import (
	"net/http"

	"github.com/google/uuid"
)

const (
	SessionIDCookie = "session_id"
)

func SessionID(r *http.Request) (uuid.UUID, bool) {
	c, err := r.Cookie(SessionIDCookie)
	if err != nil {
		return uuid.Nil, false
	}

	id, _ := uuid.Parse(c.Value)
	return id, true
}
