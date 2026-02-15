package httpapi

import (
	"context"
	"log"
	"messenger/internal/contextx"
	chatpb "messenger/proto/chats"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 30 * time.Second,
	CheckOrigin:      func(r *http.Request) bool { return true },
}

func (g *Gateway) serveWS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatId, err := requestChatId(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user := contextx.UserWithCtx(r.Context())

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("ws upgrade error:", err)
			return
		}
		defer conn.Close()

		stream, err := g.chatClient.Connect(context.Background())
		if err != nil {
			log.Println("grpc connect error:", err)
			return
		}

		err = stream.Send(&chatpb.ChatMessage{
			ChatId:   proto.String(chatId.String()),
			UserId:   user.UserId,
			Username: user.Username,
		})
		if err != nil {
			log.Println("grpc register error:", err)
			return
		}

		go func() {
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					stream.CloseSend()
					return
				}

				stream.Send(&chatpb.ChatMessage{
					ChatId:    proto.String(chatId.String()),
					UserId:    user.UserId,
					Text:      proto.String(string(message)),
					Timestamp: timestamppb.Now(),
					Username:  user.Username,
				})
			}
		}()

		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Println("grpc recv error:", err)
				return
			}

			data, _ := protojson.Marshal(msg)
			err = conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				return
			}
		}
	}
}
