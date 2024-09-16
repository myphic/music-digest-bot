package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"log"
	"music-digest-bot/internal/testhelpers"
	"testing"
)

var yaMusic = "yandexmusic"

type SourcesRepoSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repository  *SourcesRepositoryImpl
	ctx         context.Context
}

func (suite *SourcesRepoSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer

	pool, err := pgxpool.New(suite.ctx, pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	db := stdlib.OpenDBFromPool(pool)

	err = goose.Up(db, "./../migrations")
	if err != nil {
		log.Fatal(err)
	}
	sourcesRepo := NewSourcesRepository(pool)
	suite.repository = sourcesRepo
}

func (suite *SourcesRepoSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *SourcesRepoSuite) TestCreateSources() {
	t := suite.T()
	source, err := suite.repository.Create(suite.ctx, SourceModel{Name: yaMusic})
	require.NoError(t, err)
	require.NotNil(t, source)
	require.Equal(t, yaMusic, source.Name)
}

func TestSourcesRepoTestSuite(t *testing.T) {
	suite.Run(t, new(SourcesRepoSuite))
}
