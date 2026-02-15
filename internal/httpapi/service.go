package httpapi

import (
	authpb "messenger/proto/auth"
	chatpb "messenger/proto/chats"
	"net/http"

	"github.com/gorilla/mux"
)

type Gateway struct {
	authClient authpb.AuthServiceClient
	chatClient chatpb.ChatServiceClient
}

func NewGateway(
	authClient authpb.AuthServiceClient,
	chatClient chatpb.ChatServiceClient,
) *Gateway {
	return &Gateway{
		authClient: authClient,
		chatClient: chatClient,
	}
}

func (g *Gateway) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/register", g.registerHandler())

	mux.Handle("/login", g.loginHandler())

	mux.Handle("/logout", g.logoutHandler())

	logged := mux.PathPrefix("/logged").Subrouter()
	g.RegisterChatRoutes(logged)

	mux.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))
}

func (g *Gateway) RegisterChatRoutes(mux *mux.Router) {
	mux.Use(func(next http.Handler) http.Handler {
		return g.withAuth(next)
	})

	mux.HandleFunc("/ws", g.serveWS())

	mux.HandleFunc("/history", g.history())

	mux.HandleFunc("/chats", g.chats())

	mux.HandleFunc("/direct", g.startDirect())

	mux.HandleFunc("/group", g.startGroup())
}
