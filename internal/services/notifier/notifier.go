package notifier

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"music-digest-bot/internal/db/repository"
	"strings"
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
	m := make(map[string][]repository.DigestModel)
	for i := range digests {
		m[digests[i].Genre] = append(m[digests[i].Genre], digests[i])
		err = n.digest.MarkAsPosted(ctx, digests[i])
		if err != nil {
			return err
		}
	}
	err = n.sendArticle(m)
	if err != nil {
		return err
	}
	return nil
}

func (n *Notifier) sendArticle(m map[string][]repository.DigestModel) error {
	var sb strings.Builder
	for k, digest := range m {
		var tracks strings.Builder
		for i := range digest {
			tracks.WriteString(fmt.Sprintf("<b>Трек: %s</b>\n<b>Описание: %s </b> \n", digest[i].Title, digest[i].Description))
		}
		sb.WriteString(fmt.Sprintf("<b>"+
			"<b>Жанр: %s</b>\n"+
			"<b>%s</b>"+
			"</b>\n", k, tracks.String()))
	}

	msg := tgbotapi.NewMessage(n.channelID, sb.String())
	msg.ParseMode = "HTML"
	_, err := n.bot.Send(msg)

	if err != nil {
		return err
	}

	return nil
}
