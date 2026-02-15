package main

import (
	"log"
	"net/http"

	"messenger/internal/httpapi"
	authpb "messenger/proto/auth"
	chatpb "messenger/proto/chats"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"auth-service:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	authClient := authpb.NewAuthServiceClient(conn)
	log.Println("Created AuthServiceClient on 'auth-service:50052'")

	conn, err = grpc.NewClient(
		"chat-service:50053",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	chatClient := chatpb.NewChatServiceClient(conn)
	log.Println("Created ChatServiceClient on 'chat-service:50053'")

	mux := mux.NewRouter()
	service := httpapi.NewGateway(
		authClient,
		chatClient,
	)
	service.RegisterRoutes(mux)

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
