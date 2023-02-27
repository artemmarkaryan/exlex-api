package telegram

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	Token string
}

const key = "tg_bot"
const reportChatID = -853602207

func MakeBot(ctx context.Context, cfg Config) (context.Context, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("error creating bot %w", err)
	}

	bot.Debug = true

	log.Printf("telegram bot: authorized on account %s", bot.Self.UserName)

	// u := tgbotapi.NewUpdate(0)
	// u.Timeout = 120
	//
	// updates := bot.GetUpdatesChan(u)
	//
	// go func() {
	// 	for update := range updates {
	// 		if update.Message != nil { // If we got a message
	// 			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	//
	// 			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	// 			msg.ReplyToMessageID = update.Message.MessageID
	//
	// 			bot.Send(msg)
	// 		}
	// 	}
	// }()

	return context.WithValue(ctx, key, bot), nil
}

func Propagate(ctx context.Context, b *tgbotapi.BotAPI) context.Context {
	return context.WithValue(ctx, key, b)
}

func FromContext(ctx context.Context) *tgbotapi.BotAPI {
	b, ok := ctx.Value(key).(*tgbotapi.BotAPI)
	if !ok {
		panic("no bot in context")
	}

	return b
}

func Report(ctx context.Context, text string) {
	bot := FromContext(ctx)

	msg := tgbotapi.NewMessage(reportChatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("error sending report: %v", err)
	}
}
