package notifictions

import (
	"fmt"
	"github.com/slack-go/slack"
)

type SlackWebhookNotifier struct {
	Config SlackConfig
}

func (notifier *SlackWebhookNotifier) Notify(message string) error {
	if notifier.Config.WebhookURL == "" || notifier.Config.WebhookChannel == "" {
		return fmt.Errorf("WebhookURL and WebhookChannel must be specified")
	}

	slackMsg := &slack.WebhookMessage{
		Text:      message,
		Channel:   notifier.Config.WebhookChannel,
		Username:  notifier.Config.WebhookUserName,
		IconEmoji: ":mage:",
	}

	return slack.PostWebhook(notifier.Config.WebhookURL, slackMsg)
}
