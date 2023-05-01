package notifictions

import (
	"magecomm/config_manager"
)

type SlackConfig struct {
	Enabled         bool
	WebhookURL      string
	WebhookChannel  string
	WebhookUserName string
	AppToken        string
	Token           string
	ChannelID       string
}

var DefaultSlackConfig = SlackConfig{
	Enabled:         false,
	WebhookURL:      "",
	WebhookChannel:  "",
	WebhookUserName: "",
	AppToken:        "",
	Token:           "",
	ChannelID:       "",
}

var DefaultSlackNotifier Notifier

func NewSlackConfig() SlackConfig {
	return SlackConfig{
		Enabled:         config_manager.GetBoolValue(config_manager.ConfigSlackEnabled),
		WebhookURL:      config_manager.GetValue(config_manager.ConfigSlackWebhookUrl),
		WebhookChannel:  config_manager.GetValue(config_manager.ConfigSlackWebhookChannel),
		WebhookUserName: config_manager.GetValue(config_manager.ConfigSlackWebhookUserName),
		AppToken:        config_manager.GetValue(config_manager.ConfigSlackAppToken),
		Token:           config_manager.GetValue(config_manager.ConfigSlackAppChannel),
		ChannelID:       config_manager.GetValue(config_manager.ConfigSlackAppUserName),
	}
}

func NewSlackNotifier(config SlackConfig) Notifier {
	if config.Enabled && config.WebhookURL != "" {
		return &SlackWebhookNotifier{Config: config}
	}
	return &SlackAppNotifier{Config: config}
}

func init() {
	DefaultSlackConfig = NewSlackConfig()
	DefaultSlackNotifier = NewSlackNotifier(DefaultSlackConfig)
}
