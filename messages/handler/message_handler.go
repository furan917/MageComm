package handler

import (
	"fmt"
	"magecomm/logger"
	"magecomm/magerun"
)

const MessageRetryLimit = 5

func HandleReceivedMessage(queueName string, messageBody string) error {
	logger.Debugf("Handling message from queue:", queueName)

	switch queueName {
	case "magerun":
		magerun.HandleMagerunCommand(messageBody)
	case "deploy":
		logger.Infof("Deploying...")
	default:
		return fmt.Errorf("no known message handler for queue: %s", queueName)
	}
	return nil
}
