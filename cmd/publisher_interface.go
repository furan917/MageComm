package cmd

import (
	"fmt"
	publisher2 "magecomm/messages/publisher"
)

type MessagePublisher interface {
	Publish(messageBody string, queue string, addCorrelationID string) (string, error)
}

var publisher MessagePublisher = &defaultMessagePublisher{}

type defaultMessagePublisher struct{}

func (d *defaultMessagePublisher) Publish(messageBody string, queue string, addCorrelationID string) (string, error) {
	publisher, err := publisher2.MapPublisherToEngine()
	if err != nil {
		return "", fmt.Errorf("failed to map publisher to engine: %v", err)
	}

	correlationID, err := publisher.PublishMessage(messageBody, queue, addCorrelationID)
	if err != nil {
		return "", fmt.Errorf("failed to publish message: %v", err)
	}
	if correlationID != "" {
		return publisher2.HandleCorrelationID(publisher, correlationID, queue)
	}

	return "", nil
}

func SetMessagePublisher(p MessagePublisher) {
	publisher = p
}
