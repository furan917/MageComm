package handler

import (
	"fmt"
	"magecomm/config_manager"
	"magecomm/logger"
	"magecomm/magerun"
	"magecomm/messages/publisher"
	"magecomm/messages/queues"
	"magecomm/notifictions"
	"strings"
)

type MagerunHandler struct {
	Publisher publisher.MessagePublisher
}

var MageRunQueue = "magerun"

func (handler *MagerunHandler) ProcessMessage(messageBody string, correlationID string) error {

	output, err := magerun.HandleMagerunCommand(messageBody)
	if err != nil {
		output = output + err.Error()
	}

	if config_manager.GetBoolValue(config_manager.ConfigSlackEnabled) && !config_manager.GetBoolValue(config_manager.ConfigSlackDisableOutputNotifications) {
		logger.Infof("Slack notification is enabled, sending output notification")
		notifier := notifictions.DefaultSlackNotifier
		outputMessage := fmt.Sprintf(
			" Command: '%v' on environment: '%s' ran with output: \n %s",
			strings.Join(strings.Fields(messageBody), " "),
			config_manager.GetValue(config_manager.CommandConfigEnvironment),
			output)
		err := notifier.Notify(fmt.Sprintf(outputMessage))
		if err != nil {
			logger.Warnf("Failed to send slack output notification: %v\n", err)
		}
	}

	// Publish the output to the RMQ/SQS output queue
	publisherClass, err := publisher.MapPublisherToEngine()
	if err != nil {
		logger.Warnf("Error publishing message to RMQ/SQS queue: %s", err)
	}
	logger.Debugf("Publishing output to queue: %s with correlation ID: %s", queues.MapQueueToOutputQueue(MageRunQueue), correlationID)
	_, err = publisherClass.Publish(output, queues.MapQueueToOutputQueue(MageRunQueue), correlationID)
	if err != nil {
		return err
	}

	return nil
}
