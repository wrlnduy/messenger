package httpapi

import (
	"encoding/json"
	"log"
	"messenger/internal/contextx"
	"messenger/internal/cookies"
	authpb "messenger/proto/auth"
	chatpb "messenger/proto/chats"
	"net/http"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	sessionTime = time.Hour * 24 * 30
)

func (g *Gateway) registerHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := g.authClient.Register(
			r.Context(),
			&authpb.RegisterRequest{
				Username: proto.String(req.Username),
				Password: proto.String(req.Password),
			},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = g.chatClient.AddToGroup(
			r.Context(),
			&chatpb.AddToGroupRequest{
				UserId: resp.UserId,
			},
		)

		w.WriteHeader(http.StatusCreated)
	})
}

func (g *Gateway) loginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp, err := g.authClient.Login(
			r.Context(),
			&authpb.LoginRequest{
				Username: proto.String(req.Username),
				Password: proto.String(req.Password),
			},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     cookies.SessionIDCookie,
			Value:    resp.GetSessionId(),
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(sessionTime),
		})

		w.WriteHeader(http.StatusOK)
	})
}

func (g *Gateway) logoutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId, ok := cookies.SessionID(r)
		if !ok {
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     cookies.SessionIDCookie,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   0,
		})

		_, err := g.authClient.Logout(
			r.Context(),
			&authpb.LogoutRequest{
				SessionId: proto.String(sessionId.String()),
			},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("SessionId %q has been ended\n", sessionId)

		w.WriteHeader(http.StatusOK)
	})
}

// func (g *Gateway) postMessage() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req struct {
// 			Text string `json:"text"`
// 		}
// 		err := json.NewDecoder(r.Body).Decode(&req)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		chatId, err := requestChatId(r)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		_, err = g.chatClient.PostMessage(
// 			r.Context(),
// 			&chatpb.PostMessageRequest{
// 				ChatId: proto.String(chatId.String()),
// 				Text:   proto.String(req.Text),
// 				UserId: contextx.UserWithCtx(r.Context()).UserId,
// 			},
// 		)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		w.WriteHeader(http.StatusNoContent)
// 	}
// }

func (g *Gateway) history() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatId, err := requestChatId(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hist, err := g.chatClient.GetHistory(
			r.Context(),
			&chatpb.HistoryRequest{
				ChatId: proto.String(chatId.String()),
				UserId: contextx.UserWithCtx(r.Context()).UserId,
			},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, _ := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(hist)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func (g *Gateway) chats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chats, err := g.chatClient.GetChats(
			r.Context(),
			&chatpb.GetChatsRequest{
				UserId: contextx.UserWithCtx(r.Context()).UserId,
			},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		data, _ := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(chats)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func (g *Gateway) startDirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		chat, err := g.chatClient.StartDirect(
			r.Context(),
			&chatpb.StartDirectRequest{
				Buddy:  proto.String(req.Username),
				UserId: contextx.UserWithCtx(r.Context()).UserId,
			},
		)

		data, _ := protojson.Marshal(chat)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func (g *Gateway) startGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Title     string   `json:"title"`
			Usernames []string `json:"usernames"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		chat, err := g.chatClient.StartGroup(
			r.Context(),
			&chatpb.StartGroupRequest{
				Title:     proto.String(req.Title),
				Usernames: req.Usernames,
				CreatorId: contextx.UserWithCtx(r.Context()).UserId,
			},
		)

		data, _ := protojson.Marshal(chat)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}
