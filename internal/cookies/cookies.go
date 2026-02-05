package cookies

import "net/http"

const (
	SessionIDCookie = "session_id"
)

func SessionID(r *http.Request) (string, bool) {
	c, err := r.Cookie(SessionIDCookie)
	if err != nil {
		return "", false
	}
	return c.Value, true
}
