package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/require"
	"music-digest-bot/internal/testhelpers"
	"testing"
)

var yaMusic = "yandexmusic"

func TestSourcesRepository(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := testhelpers.CreatePostgresContainer(ctx)
	require.NoError(t, err)

	dbUrl := pgContainer.ConnectionString

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
