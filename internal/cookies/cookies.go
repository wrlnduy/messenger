package cookies

import "net/http"

const UserIDCookie = "user_id"

func UserID(r *http.Request) string {
	c, _ := r.Cookie(UserIDCookie)
	return c.Value
}
