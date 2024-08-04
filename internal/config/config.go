package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseUrl       string `mapstructure:"DATABASE_URL"`
	TelegramBotToken  string `mapstructure:"TELEGRAM_BOT_TOKEN"`
	TelegramChannelID int64  `mapstructure:"TELEGRAM_CHANNEL_ID"`
	YandexMusicToken  string `mapstructure:"YANDEX_MUSIC_TOKEN"`
}

func FromEnv(path string) (*Config, error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName("app")
	v.SetConfigType("env")
	v.SetDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres")
	v.AutomaticEnv()
	err := v.ReadInConfig()

	if err != nil {
		return nil, err
	}

	cfg := Config{}
	err = v.Unmarshal(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
