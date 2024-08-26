package notifier

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"music-digest-bot/internal/db/repository"
	"time"
)

type DigestProvider interface {
	AllNotPosted(ctx context.Context) ([]repository.DigestModel, error)
	MarkAsPosted(ctx context.Context, article repository.DigestModel) error
}

type Notifier struct {
	digest       *repository.DigestRepositoryImpl
	sendInterval time.Duration
	channelID    int64
	bot          *tgbotapi.BotAPI
}

func New(digest *repository.DigestRepositoryImpl, sendInterval time.Duration, channelID int64, bot *tgbotapi.BotAPI) *Notifier {
	return &Notifier{
		digest:       digest,
		sendInterval: sendInterval,
		channelID:    channelID,
		bot:          bot,
	}
}

func (n *Notifier) Start(ctx context.Context) error {
	ticker := time.NewTicker(n.sendInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := n.SelectAndSendArticle(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (n *Notifier) SelectAndSendArticle(ctx context.Context) error {
	digests, err := n.digest.AllNotPosted(ctx)
	if err != nil {
		return err
	}

	for i := range digests {
		err = n.sendArticle(digests[i])
		if err != nil {
			return err
		}
		err = n.digest.MarkAsPosted(ctx, digests[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *Notifier) sendArticle(digest repository.DigestModel) error {
	const msgFormat = "*%s\n\n%s\n\n%s\n\n%s"

	msg := tgbotapi.NewMessage(n.channelID, fmt.Sprintf(
		msgFormat,
		digest.Title,
		digest.Genre,
		digest.PublishedAt,
		digest.Description,
	))

	_, err := n.bot.Send(msg)

	if err != nil {
		return err
	}

	return nil
}
