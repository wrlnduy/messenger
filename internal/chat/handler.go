package chat

import (
	"encoding/json"
	"net/http"
	"time"

	"messenger/internal/storage"
	"messenger/internal/ws"
)

type request struct {
	Text string `json:"text"`
}

type response struct {
	UserID    string `json:"user_id"`
	Text      string `json:"text"`
	Timestamp int64  `json:"timestamp"`
}

func PostMessage(hub *ws.Hub, store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}

		msg := storage.Message{
			UserID:    "anon",
			Text:      req.Text,
			CreatedAt: time.Now(),
		}

		store.Save(msg)

		resp := response{
			UserID:    msg.UserID,
			Text:      msg.Text,
			Timestamp: msg.CreatedAt.Unix(),
		}

		data, _ := json.Marshal(resp)
		hub.Broadcast(data)

		w.WriteHeader(http.StatusNoContent)
	}
}
