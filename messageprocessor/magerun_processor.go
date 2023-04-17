package messageprocessor

import (
	"magecomm/logger"
	"magecomm/magerun"
	"magecomm/messages/publisher"
	"magecomm/messages/queues"
)

type MagerunProcessor struct {
	Publisher publisher.MessagePublisher
}

func (process *MagerunProcessor) ProcessMessage(messageBody string, correlationID string) error {

	output, err := magerun.HandleMagerunCommand(messageBody)
	if err != nil {
		return err
	}

	// Publish the output to the RMQ/SQS output queue
	publisherClass, err := publisher.MapPublisherToEngine()
	if err != nil {
		logger.Warnf("Error publishing message to RMQ/SQS queue:", err)
	}
	_, err = publisherClass.Publish(output, queues.MapQueueToOutputQueue(magerun.CommandMageRun), correlationID)
	if err != nil {
		return err
	}

	return nil
}
