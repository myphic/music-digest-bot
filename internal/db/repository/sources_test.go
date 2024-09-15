package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
	"time"
)

var yaMusic = "yandexmusic"

func TestSourcesRepository(t *testing.T) {
	ctx := context.Background()

	dbName := "music_digest_bot"
	dbUser := "postgres"
	dbPassword := "postgres"

	postgresContainer, err := postgres.Run(ctx,
		"docker.io/postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	dbUrl, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	pool, err := pgxpool.New(ctx, dbUrl)
	require.NoError(t, err)
	db := stdlib.OpenDBFromPool(pool)

	err = goose.Up(db, "./../migrations")
	require.NoError(t, err)
	sourcesRepo := NewSourcesRepository(pool)

	source, err := sourcesRepo.Create(ctx, SourceModel{Name: yaMusic})
	require.NoError(t, err)
	require.NotNil(t, source)
	require.Equal(t, yaMusic, source.Name)
}
