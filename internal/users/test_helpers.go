package users

import (
	"context"
	"database/sql"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupStore(t *testing.T) (Store, func()) {
	ctx := context.Background()

	pgC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16",
			Env:          map[string]string{"POSTGRES_PASSWORD": "test", "POSTGRES_USER": "test", "POSTGRES_DB": "messenger_test"},
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor:   wait.ForListeningPort("5432/tcp"),
		},
		Started: true,
	})
	require.NoError(t, err)

	host, _ := pgC.Host(ctx)
	port, _ := pgC.MappedPort(ctx, "5432")
	dsn := "postgres://test:test@" + host + ":" + port.Port() + "/messenger_test?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	cleanup := func() {
		db.Close()
		pgC.Terminate(ctx)
	}

	store, err := NewPostgresStore(db)
	if err != nil {
		cleanup()
	}
	require.NoError(t, err)

	return store, cleanup
}

func setupCache(t *testing.T, store Store) (*UserCache, func()) {
	ctx := context.Background()

	rdC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:8.4-alpine",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForListeningPort("6379/tcp"),
		},
		Started: true,
	})
	require.NoError(t, err)

	endpoint, err := rdC.Endpoint(ctx, "")
	require.NoError(t, err)

	rdb := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})
	require.NoError(t, rdb.Ping(ctx).Err())

	cache, err := NewUserCache(rdb, store)
	require.NoError(t, err)

	cleanup := func() {
		rdb.Close()
		rdC.Terminate(ctx)
	}

	return cache, cleanup
}

func SetupService(t *testing.T) (*Service, func()) {
	store, sCleanup := setupStore(t)
	cache, cCleanup := setupCache(t, store)

	cleanup := func() {
		sCleanup()
		cCleanup()
	}

	svc := NewService(store, cache)
	return svc, cleanup
}
