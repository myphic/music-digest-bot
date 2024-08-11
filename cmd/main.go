package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"music-digest-bot/internal/config"
	"music-digest-bot/internal/db/repository"
	"music-digest-bot/internal/services"
	"music-digest-bot/internal/services/yandexmusic"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	cfg, err := config.FromEnv(".")
	if err != nil {
		logger.Error("error getting config", err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	yaMusic := yandexmusic.NewYandexFetcher(cfg.YandexMusicToken)

	conn, err := pgx.Connect(ctx, cfg.DatabaseUrl)
	if err != nil {
		logger.Error("database connection error: ", err)
	}

	defer conn.Close(ctx)
	sourcesRepo := repository.NewSourcesRepository(conn)
	digestRepo := repository.NewDigestRepository(conn)
	fetcher := services.New(sourcesRepo, digestRepo, yaMusic)
	err = fetcher.Fetch(ctx)
	if err != nil {
		logger.Error("error fetching sources: ", err)
	}
}
