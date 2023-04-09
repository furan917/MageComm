package cmd

import (
	"github.com/spf13/cobra"
	"magecomm/config_manager"
	"magecomm/logger"
	"magecomm/messages"
	"strings"
)

var ListenCmd = &cobra.Command{
	Use:   "listen [queue1] [queue2] ...",
	Short: "Listen for messages from specified queues, fallback to ENV LISTENERS, use -e or ENV LISTENER_ENGINE to specify engine (sqs|rmq), default sqs",
	Run: func(cmd *cobra.Command, args []string) {
		queueNames := args
		if len(queueNames) == 0 {
			queuesFromEnv := config_manager.GetValue("listeners")
			if queuesFromEnv == "" {
				logger.Fatal("No queues specified")
				return
			}
			queueNames = strings.Split(queuesFromEnv, ",")
		}

		messages.MapListenerToEngine(queueNames)
	},
}
