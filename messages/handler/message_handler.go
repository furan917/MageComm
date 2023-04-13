package handler

import (
	"fmt"
	"magecomm/logger"
	"magecomm/magerun"
)

const MessageRetryLimit = 5

func HandleReceivedMessage(messageBody string, queueName string, correlationID string) error {
	logger.Debugf("Handling message from queue:", queueName)

	switch queueName {
	case "magerun":
		magerun.HandleMagerunCommand(messageBody, correlationID)
	case "deploy":
		logger.Infof("Deploying...")
	default:
		return fmt.Errorf("no known message handler for queue: %s", queueName)
	}
	return nil
}
