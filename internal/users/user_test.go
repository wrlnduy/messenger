package users

import (
	"context"
	"testing"

	_ "github.com/lib/pq"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStore_CreateAndFindUser(t *testing.T) {
	ctx := context.Background()
	store, cleanup := setupStore(t)
	defer cleanup()

	userId := uuid.New()
	username := "aba"
	passwordHash := "caba"

	err := store.CreateUser(ctx, userId, username, passwordHash)
	require.NoError(t, err)

	u, err := store.FindByID(ctx, userId)
	require.NoError(t, err)
	require.Equal(t, username, *u.Username)
	require.Equal(t, userId.String(), *u.UserId)

	u, err = store.FindByUsername(ctx, username)
	require.NoError(t, err)
	require.Equal(t, username, *u.Username)
	require.Equal(t, userId.String(), *u.UserId)

	_, err = store.FindByUsername(ctx, "abacaba")
	require.Error(t, err)
}

func TestService_CreateAndFindUser(t *testing.T) {
	ctx := context.Background()
	svc, cleanup := SetupService(t)
	defer cleanup()

	username := "aba"
	passwordHash := "caba"

	userId, err := svc.CreateUser(ctx, username, passwordHash)
	require.NoError(t, err)

	u, err := svc.FindByID(ctx, userId)
	require.NoError(t, err)
	require.Equal(t, username, *u.Username)
	require.Equal(t, userId.String(), *u.UserId)

	u, err = svc.FindByUsername(ctx, username)
	require.NoError(t, err)
	require.Equal(t, username, *u.Username)
	require.Equal(t, userId.String(), *u.UserId)

	_, err = svc.FindByID(ctx, uuid.New())
	require.Error(t, err)

	_, err = svc.FindByUsername(ctx, "abacaba")
	require.Error(t, err)
}

func TestService_GetMapping(t *testing.T) {
	ctx := context.Background()
	svc, cleanup := SetupService(t)
	defer cleanup()

	username := "aba"
	passwordHash := "caba"

	userId, err := svc.CreateUser(ctx, username, passwordHash)
	require.NoError(t, err)

	userIDs := []uuid.UUID{userId}
	mapping, err := svc.GetMapping(ctx, userIDs)
	require.NoError(t, err)
	require.Equal(t, username, mapping[userId])

	userIDs = append(userIDs, uuid.New())
	mapping, err = svc.GetMapping(ctx, userIDs)
	require.Error(t, err)
}
