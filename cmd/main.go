package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"music-digest-bot/internal/config"
	"music-digest-bot/internal/db/repository"
	"music-digest-bot/internal/services/yandexmusic"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.FromEnv(".")
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	yaMusic := yandexmusic.Fetch{}
	fmt.Println(cfg, yaMusic)
	conn, err := pgx.Connect(ctx, cfg.DatabaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)
	repo := repository.NewSourcesRepository(conn)
	sources, err := repo.Sources(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sources)
	//yaMusic.Fetch(ctx, cfg.YandexMusicToken)
}
