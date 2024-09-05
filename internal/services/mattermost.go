package services

import (
	mmSDK "github.com/mattermost/mattermost-server/v6/model"
)

type MattermostNotifier struct {
	bot       *mmSDK.Client4
	channelID string
}

func (mm *MattermostNotifier) Notify(message string) error {
	post := &mmSDK.Post{
		ChannelId: mm.channelID,
		Message:   message,
	}
	_, _, err := mm.bot.CreatePost(post)
	if err != nil {
		return err
	}

	return nil
}
