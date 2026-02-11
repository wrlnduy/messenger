package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"messenger/internal/auth"
	"messenger/internal/cache"
	"messenger/internal/messages"
	"messenger/internal/ws"
	messenger "messenger/proto"

	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type request struct {
	Text string `json:"text"`
}

func PostMessage(hub *ws.Hub, store messages.Store, cache *cache.UserCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		user := auth.UserWithCtx(r.Context())
		msg := &messenger.ChatMessage{
			MessageId: proto.String(uuid.NewString()),
			UserId:    proto.String(*user.UserId),
			Text:      proto.String(req.Text),
			Timestamp: proto.Int64(time.Now().Unix()),
		}

		err = store.Save(r.Context(), msg)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userId, _ := uuid.Parse(*user.UserId)
		username, err := cache.GetUsername(r.Context(), userId)
		msg.Username = proto.String(username)

		data, _ := protojson.Marshal(msg)
		hub.Broadcast(data)

		w.WriteHeader(http.StatusNoContent)
	}
}

func History(store messages.Store, cache *cache.UserCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hist, err := store.History(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = FillMapping(r.Context(), hist, cache)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, _ := protojson.Marshal(hist)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}
