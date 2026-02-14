package chat

import (
	"database/sql"
	"encoding/json"
	"errors"
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

func PostMessage(manager *ws.HubManager, store messages.Store, cache *cache.UserCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Text string `json:"text"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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

		data, _ := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(hist)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func Chats(store chats.Store, cache *cache.UserCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := auth.UserIdWithCtx(r.Context())

		chats, err := store.GetUserChats(r.Context(), userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, chat := range chats.Chats {
			if chat.Type.Number() != messenger.ChatType_DIRECT.Number() {
				continue
			}

			chatId, _ := uuid.Parse(*chat.ChatId)
			buddyId, err := store.GetDirectBuddyId(r.Context(), chatId, userId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			username, err := cache.GetUsername(r.Context(), buddyId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			chat.Title = proto.String(username)
		}

		data, _ := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(chats)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func GetCreateDirect(chats chats.Store, cache *cache.UserCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u1 := auth.UserIdWithCtx(r.Context())

		var req struct {
			Username string `json:"username"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u2, err := cache.GetUserId(r.Context(), req.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		chat, err := chats.GetDirect(r.Context(), u1, u2)
		if errors.Is(err, sql.ErrNoRows) {
			chat, err = chats.CreateDirect(r.Context(), u1, u2)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, _ := protojson.Marshal(chat)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func CreateGroup(chats chats.Store, cache *cache.UserCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creator := auth.UserIdWithCtx(r.Context())

		var req struct {
			Title     string   `json:"title"`
			Usernames []string `json:"usernames"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		users, err := UserIDsByUsernames(r.Context(), cache, req.Usernames)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		chat, err := chats.CreateGroup(r.Context(), creator, req.Title, users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, _ := protojson.Marshal(chat)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}
