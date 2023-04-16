package cmd

import (
	"github.com/spf13/cobra"
	"magecomm/config_manager"
	"magecomm/logger"
	"magecomm/messages/listener"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var ListenCmd = &cobra.Command{
	Use:   "listen [queue1] [queue2] ...",
	Short: "Listen for messages from specified queues, fallback to ENV LISTENERS, use -e or ENV LISTENER_ENGINE to specify engine (sqs|rmq), default sqs",
	Run: func(cmd *cobra.Command, args []string) {
		queueNames := args
		if len(queueNames) == 0 {
			queuesFromEnv := config_manager.GetValue(config_manager.CommandConfigListeners)
			if queuesFromEnv == "" {
				logger.Fatal("No queues specified")
				return
			}
			queueNames = strings.Split(queuesFromEnv, ",")
		}

		listener, err := listener.MapListenerToEngine()
		if err != nil {
			logger.Fatal(err)
			return
		}

		// Create a channel to handle program termination or interruption signals so we can kill any connections if needed
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		go listener.ListenToService(queueNames)
		<-sigChan
		listener.Close()
	},
}
