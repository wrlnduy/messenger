package usergrpc

import (
	"context"
	"messenger/internal/users"
	userpb "messenger/proto/users"
	"net"
	"testing"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func startGRPCServer(t *testing.T, svc *users.Service) (userpb.UsersServiceClient, func()) {
	lis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	grpcServer := grpc.NewServer()
	userpb.RegisterUsersServiceServer(grpcServer, NewServer(svc))

	go grpcServer.Serve(lis)

	cleanup := func() {
		grpcServer.Stop()
		lis.Close()
	}

	conn, err := grpc.NewClient(
		lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	client := userpb.NewUsersServiceClient(conn)

	return client, func() {
		conn.Close()
		cleanup()
	}
}

func TestGRPC_CreateAndFindUser(t *testing.T) {
	ctx := context.Background()
	svc, cleanup := users.SetupService(t)
	defer cleanup()

	client, grpcCleanup := startGRPCServer(t, svc)
	defer grpcCleanup()

	username := "aba"
	passwordHash := "caba"

	createResp, err := client.CreateUser(ctx, &userpb.CreateUserRequest{
		Username:     proto.String(username),
		PasswordHash: proto.String(passwordHash),
	})
	require.NoError(t, err)
	require.NotNil(t, createResp)

	userByName, err := client.FindByUsername(ctx, &userpb.FindByUsernameRequest{
		Username: proto.String(username),
	})
	require.NoError(t, err)
	require.Equal(t, username, *userByName.Username)
	userId, _ := uuid.Parse(*userByName.UserId)

	userByID, err := client.FindByID(ctx, &userpb.FindByIDRequest{
		UserId: proto.String(userId.String()),
	})
	require.NoError(t, err)
	require.Equal(t, username, *userByID.Username)

	_, err = client.FindByID(ctx, &userpb.FindByIDRequest{
		UserId: proto.String(uuid.NewString()),
	})
	require.Error(t, err)
}
