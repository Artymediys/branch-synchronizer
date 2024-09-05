package services

import (
	"errors"

	"branch-synchronizer/env"

	mmSDK "github.com/mattermost/mattermost-server/v6/model"
	"github.com/mymmrac/telego"
)

type Notifier interface {
	Notify(message string) error
}

type NotifierClient struct {
	notifiers []Notifier
}

func (nc *NotifierClient) SendNotification(message string) error {
	for _, notifier := range nc.notifiers {
		err := notifier.Notify(message)
		if err != nil {
			return errors.New("error sending notification: " + err.Error())
		}
	}

	return nil
}

func NewNotifierClient(config *env.Config) (*NotifierClient, error) {
	var notifiers []Notifier

	if config.Telegram.Enabled {
		tgBot, err := telego.NewBot(config.Telegram.BotToken)
		if err != nil {
			return nil, errors.New("error creating Telegram bot client: " + err.Error())
		}

		notifiers = append(notifiers, &TelegramNotifier{bot: tgBot, channelID: config.Telegram.ChannelID})
	}

	if config.Mattermost.Enabled {
		mmBot := mmSDK.NewAPIv4Client(config.Mattermost.Url)
		mmBot.SetOAuthToken(config.Mattermost.BotToken)

		notifiers = append(notifiers, &MattermostNotifier{bot: mmBot, channelID: config.Mattermost.ChannelID})
	}

	if len(notifiers) == 0 {
		return nil, errors.New("no notifiers configured")
	}

	return &NotifierClient{notifiers: notifiers}, nil
}
