package services

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type TelegramNotifier struct {
	bot       *telego.Bot
	channelID int64
}

func (tg *TelegramNotifier) Notify(message string) error {
	msg := tu.Message(
		tu.ID(tg.channelID),
		message,
	).WithParseMode("Markdown")
	_, err := tg.bot.SendMessage(msg)

	return err
}
