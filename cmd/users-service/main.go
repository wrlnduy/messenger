package main

import (
	"log"
	"messenger/internal/db"
	"messenger/internal/users"
	"net"

	usergrpc "messenger/internal/users/grpc"
	userspb "messenger/proto/users"

	"google.golang.org/grpc"
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

	usersStore, err := users.NewPostgresStore(dbs)
	if err != nil {
		log.Fatal(err)
	}

	userCache, err := users.NewUserCache(rdb, usersStore)
	if err != nil {
		log.Fatal(err)
	}

	service := users.NewService(usersStore, userCache)

	grpcServer := grpc.NewServer()
	userspb.RegisterUsersServiceServer(grpcServer, usergrpc.NewServer(service))

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()
	log.Println("UsersService running on :50051")

	grpcServer.Serve(lis)
}
