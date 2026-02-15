package main

import (
	"log"
	"messenger/internal/auth"
	authgrpc "messenger/internal/auth/grpc"
	"messenger/internal/auth/sessions"
	"messenger/internal/db"
	authpb "messenger/proto/auth"
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

	rdb, err := db.NewRdb("REDIS_URL")
	if err != nil {
		log.Fatal(err)
	}
	defer rdb.Close()

	store, err := sessions.NewPostgresStore(dbs)
	if err != nil {
		log.Fatal(err)
	}

	cache := sessions.NewSessionsCache(rdb, store)

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

	service := auth.NewService(store, cache, usersClient)

	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, authgrpc.NewServer(service))

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()
	log.Println("AuthService running on :50052")

	grpcServer.Serve(lis)
}
