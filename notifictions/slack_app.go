package notifictions

import (
	"context"
	"fmt"
	"github.com/slack-go/slack"
)

type SlackAppNotifier struct {
	Config SlackConfig
}

func (notifier *SlackAppNotifier) Notify(message string) error {
	if notifier.Config.Token == "" || notifier.Config.AppToken == "" || notifier.Config.ChannelID == "" {
		return fmt.Errorf("token, AppToken and ChannelID must be specified")
	}

	client := slack.New(
		notifier.Config.AppToken,
		slack.OptionAppLevelToken(notifier.Config.Token),
	)

	_, _, err := client.PostMessageContext(
		context.Background(),
		notifier.Config.ChannelID,
		slack.MsgOptionText(message, false),
	)
	if err != nil {
		return fmt.Errorf("error sending Slack message: %v", err)
	}

	return nil
}
