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

func (handler *MagerunHandler) ProcessMessage(messageBody string, correlationID string) error {

	output, err := magerun.HandleMagerunCommand(messageBody)
	if err != nil {
		//ensure that any error is passed back to the caller
		output = output + err.Error()
	}

	// Publish the output to the RMQ/SQS output queue
	publisherClass, err := publisher.MapPublisherToEngine()
	if err != nil {
		logger.Warnf("Error publishing message to RMQ/SQS queue:", err)
	}
	logger.Debugf("Publishing output to queue:", queues.MapQueueToOutputQueue(magerun.CommandMageRun), "with correlation ID:", correlationID)
	_, err = publisherClass.Publish(output, queues.MapQueueToOutputQueue(magerun.CommandMageRun), correlationID)
	if err != nil {
		return err
	}

	return nil
}
