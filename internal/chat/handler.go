package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"messenger/internal/auth"
	"messenger/internal/cache"
	"messenger/internal/chats"
	"messenger/internal/messages"
	"messenger/internal/ws"
	messenger "messenger/proto"

	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type request struct {
	Text string `json:"text"`
}

func PostMessage(manager *ws.HubManager, store messages.Store, cache *cache.UserCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		hub, err := manager.GetHubForRequst(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		user := auth.UserWithCtx(r.Context())
		msg := &messenger.ChatMessage{
			MessageId: proto.String(uuid.NewString()),
			UserId:    proto.String(*user.UserId),
			ChatId:    proto.String(hub.ChatId.String()),
			Text:      proto.String(req.Text),
			Timestamp: timestamppb.New(time.Now()),
		}

		err = store.Save(r.Context(), msg)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		msg.Username = user.Username

		data, _ := protojson.Marshal(msg)
		hub.Broadcast(data)

		w.WriteHeader(http.StatusNoContent)
	}
}

func History(manager *ws.HubManager, store messages.Store, cache *cache.UserCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hub, err := manager.GetHubForRequst(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		hist, err := store.History(r.Context(), hub.ChatId)
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

func Chats(store chats.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := auth.UserIdWithCtx(r.Context())

		chats, err := store.GetUserChats(r.Context(), userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, _ := protojson.Marshal(chats)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}
