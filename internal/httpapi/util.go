package httpapi

import (
	"net/http"

	"github.com/google/uuid"
)

func requestChatId(r *http.Request) (uuid.UUID, error) {
	chatIdStr := r.URL.Query().Get("chat_id")
	return uuid.Parse(chatIdStr)
}
