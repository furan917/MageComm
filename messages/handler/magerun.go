package handler

import (
	"magecomm/logger"
	"magecomm/magerun"
	"magecomm/messages/publisher"
	"magecomm/messages/queues"
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
