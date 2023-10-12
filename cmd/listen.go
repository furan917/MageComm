package cmd

import (
	"fmt"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		queueNames := args
		if len(queueNames) == 0 {
			queuesFromConfig := config_manager.GetValue(config_manager.CommandConfigListeners)
			if queuesFromConfig == "" {
				return fmt.Errorf("no queues specified")
			}
			logger.Infof("No queues specified, using queues from Config: %s", queuesFromConfig)
			fmt.Printf("No queues specified, using queues from Config: %s", queuesFromConfig)
			queueNames = strings.Split(queuesFromConfig, ",")
		}

		//if queueNames not in allowed queues, return error
		for _, queueName := range queueNames {
			if !config_manager.IsAllowedQueue(queueName) {
				return fmt.Errorf("queue '%s' is not allowed, allowed queues: %s", queueName, config_manager.GetAllowedQueues())
			}
		}

		listener, err := listener.MapListenerToEngine()
		if err != nil {
			return fmt.Errorf("error creating listener: %s", err)
		}

		// Create a channel to handle program termination or interruption signals so we can kill any connections if needed
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		go listener.ListenToService(queueNames)
		<-sigChan
		listener.Close()

		return nil
	},
}
