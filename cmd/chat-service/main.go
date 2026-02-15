package main

import (
	"log"
	"messenger/internal/chat"
	"messenger/internal/chat/chats"
	chatgrpc "messenger/internal/chat/grpc"
	"messenger/internal/chat/members"
	"messenger/internal/chat/messages"
	storemanager "messenger/internal/chat/store_manager"
	"messenger/internal/chat/stream"
	"messenger/internal/chat/unread"
	"messenger/internal/db"
	authpb "messenger/proto/auth"
	chatpb "messenger/proto/chats"
	userpb "messenger/proto/users"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	dbs, err := db.NewDb("DATABASE_URL")
	if err != nil {
		log.Fatal(err)
	}
	defer dbs.Close()

	chats, err := chats.NewPostgresStore(dbs)
	if err != nil {
		log.Fatal(err)
	}

	hubManager := stream.NewHubManager(chats)

	messages, err := messages.NewPostgresStore(dbs)
	if err != nil {
		log.Fatal(err)
	}

	members, err := members.NewPostgresStore(dbs)
	if err != nil {
		log.Fatal(err)
	}

	unread, err := unread.NewPostgresStore(dbs)
	if err != nil {
		log.Fatal(err)
	}

	storeManager := storemanager.NewPostgresStore(
		dbs,
		messages,
		chats,
		members,
		unread,
	)

	conn, err := grpc.NewClient(
		"users-service:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	usersClient := userpb.NewUsersServiceClient(conn)
	log.Println("Created UsersServiceClient on 'users-service:50051'")

	conn, err = grpc.NewClient(
		"auth-service:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	authClient := authpb.NewAuthServiceClient(conn)
	log.Println("Created AuthServiceClient on 'auth-service:50052'")

	service := chat.NewService(
		hubManager,
		storeManager,
		messages,
		chats,
		members,
		unread,
		usersClient,
		authClient,
	)

	grpcServer := grpc.NewServer()
	chatpb.RegisterChatServiceServer(grpcServer, chatgrpc.NewServer(service))

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("ChatService running on :50053")

	grpcServer.Serve(lis)
}
