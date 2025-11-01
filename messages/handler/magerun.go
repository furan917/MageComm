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
	output = strings.TrimSpace(output)

	//if output is empty return "command finished with no out
	if output == "" {
		output = "Command finished with no output"
	}

	if config_manager.GetBoolValue(config_manager.ConfigSlackEnabled) {
		logger.Infof("Slack notification is enabled, sending output notification")
		notifier := notifictions.DefaultSlackNotifier
		commandOutputMessage := ""
		outputMessage := fmt.Sprintf(
			"Executed Command: '%v' on Environment: '%s'",
			strings.Join(strings.Fields(messageBody), " "),
			config_manager.GetValue(config_manager.CommandConfigEnvironment))

		if !config_manager.GetBoolValue(config_manager.ConfigSlackDisableOutputNotifications) {
			commandOutputMessage = fmt.Sprintf(
				"Returned with output: \n %s",
				output)
			outputMessage = outputMessage + "\n" + commandOutputMessage
		}

		err := notifier.Notify(fmt.Sprintf(outputMessage))
		if err != nil {
			logger.Warnf("Failed to send slack output notification: %v\n", err)
		}
	}

	publisherClass, err := publisher.MapPublisherToEngine()
	if err != nil {
		logger.Warnf("Error publishing message to RMQ/SQS queue: %s", err)
	}

	outputWithCommand := fmt.Sprintf("Command: %s\n\n%s", messageBody, output)
	logger.Infof("Publishing output to queue: %s with correlation ID: %s", queues.MapQueueToOutputQueue(MageRunQueue), correlationID)
	_, err = publisherClass.Publish(outputWithCommand, queues.MapQueueToOutputQueue(MageRunQueue), correlationID)
	if err != nil {
		return err
	}

	return nil
}
