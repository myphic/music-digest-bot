package main

import (
	"log"
	"music-digest-bot/internal/config"
	"music-digest-bot/internal/services/yandexmusic"
)

func main() {
	cfg, err := config.FromEnv(".")
	if err != nil {
		log.Fatal(err)
	}
	yaMusic := yandexmusic.Fetch{}

	yaMusic.Fetch(cfg.YandexMusicToken)
}
