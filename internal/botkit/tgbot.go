package botkit

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBot struct {
	api   *tgbotapi.BotAPI
	views map[string]ViewFunc
}

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

func New(api *tgbotapi.BotAPI) *TgBot {
	return &TgBot{
		api: api,
	}
}
