package main

import (
	"fmt"
	"log"
	"music-digest-bot/internal/config"
	"music-digest-bot/internal/services/yandexmusic"
)

func main() {
	cfg, err := config.FromEnv(".")
	if err != nil {
		log.Fatal(err)
	}
	ya := yandexmusic.Fetch{}
	fmt.Println(cfg.YandexMusicToken)
	ya.Fetch(cfg.YandexMusicToken)
}
