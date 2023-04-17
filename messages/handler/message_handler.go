package handler

import (
	"fmt"
	"magecomm/logger"
	"magecomm/messageprocessor"
	"magecomm/messages/publisher"
)

// todo:: Maybe move this into messageprocessor or move the messageprocessor into this package
const MessageRetryLimit = 5

func HandleReceivedMessage(messageBody string, queueName string, correlationID string) error {
	logger.Debugf("Handling message from queue:", queueName)
	var processor messageprocessor.MessageProcessor

	switch queueName {
	case "magerun":
		processor = &messageprocessor.MagerunProcessor{
			Publisher: publisher.Publisher,
		}
	case "deploy":
		logger.Infof("Deploying...")
		// assign the appropriate processor here
	default:
		return fmt.Errorf("no known message handler for queue: %s", queueName)
	}

	if processor == nil {
		return fmt.Errorf("no message processor found for queue: %s", queueName)
	}
	err := processor.ProcessMessage(messageBody, correlationID)
	if err != nil {
		return err
	}

	return nil
}
