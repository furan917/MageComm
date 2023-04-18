package publisher

import (
	"fmt"
)

type MessagePublisher interface {
	Publish(messageBody string, queueName string, addCorrelationID string) (string, error)
}

var Publisher MessagePublisher = &DefaultMessagePublisher{}

type DefaultMessagePublisher struct{}

func (d *DefaultMessagePublisher) Publish(messageBody string, queueName string, addCorrelationID string) (string, error) {
	publisher, err := MapPublisherToEngine()
	if err != nil {
		return "", fmt.Errorf("failed to map publisher to engine: %v", err)
	}

	correlationID, err := publisher.Publish(messageBody, queueName, addCorrelationID)
	if err != nil {
		return "", fmt.Errorf("failed to publish message: %v", err)
	}

	return correlationID, nil
}

func SetMessagePublisher(p MessagePublisher) {
	Publisher = p
}
