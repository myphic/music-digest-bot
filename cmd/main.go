package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"music-digest-bot/internal/config"
	"music-digest-bot/internal/db/repository"
	"music-digest-bot/internal/services"
	"music-digest-bot/internal/services/notifier"
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
		return
	}
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		logger.Error("error creating bot", err)
		return
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	yaMusic := yandexmusic.NewYandexFetcher(cfg.YandexMusicToken, logger)

	conn, err := pgx.Connect(ctx, cfg.DatabaseUrl)
	if err != nil {
		logger.Error("database connection error: ", err)
		return
	}

	defer conn.Close(ctx)
	sourcesRepo := repository.NewSourcesRepository(conn)
	digestRepo := repository.NewDigestRepository(conn)
	fetcher := services.New(sourcesRepo, digestRepo, yaMusic)
	err = fetcher.Fetch(ctx)
	if err != nil {
		logger.Error("error fetching sources: ", err)
	}
	notifier := notifier.New(digestRepo, 1000, cfg.TelegramChannelID, bot)
	err = notifier.Start(ctx)
	if err != nil {
		logger.Error("error notifying telegram channel: ", err)
	}
}
