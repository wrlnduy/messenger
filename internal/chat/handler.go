package chat

import (
	"encoding/json"
	"net/http"
	"time"

	"messenger/internal/cookies"
	"messenger/internal/storage"
	"messenger/internal/ws"
	"messenger/proto"

	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type request struct {
	Text string `json:"text"`
}

func PostMessage(hub *ws.Hub, store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}

		msg := &message.ChatMessage{
			MessageId: proto.String(uuid.NewString()),
			UserId:    proto.String(cookies.UserID(r)),
			Text:      proto.String(req.Text),
			Timestamp: proto.Int64(time.Now().Unix()),
		}

		store.Save(msg)

		data, _ := protojson.Marshal(msg)
		hub.Broadcast(data)

		w.WriteHeader(http.StatusNoContent)
	}
}
