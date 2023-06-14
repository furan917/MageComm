package handler

import (
	"fmt"
	"magecomm/logger"
	"magecomm/messages/publisher"
)

const MessageRetryLimit = 5

func HandleReceivedMessage(messageBody string, queueName string, correlationID string) error {
	logger.Debugf("Handling message from queue: %s", queueName)
	var processor MessageHandler

	switch queueName {
	case "magerun":
		processor = &MagerunHandler{
			Publisher: publisher.Publisher,
		}
	case "magerun_output":
		processor = &MagerunOutputHandler{}
	case "deploy":
		logger.Infof("Deploy not yet implemented...")
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
