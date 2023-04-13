package cmd

import (
	"magecomm/messages/mappers/publisher_mapper"
)

type MessagePublisher interface {
	Publish(messageBody string, queue string, addCorrelationID string) (string, error)
}

var publisher MessagePublisher = &defaultMessagePublisher{}

type defaultMessagePublisher struct{}

func (d *defaultMessagePublisher) Publish(messageBody string, queue string, addCorrelationID string) (string, error) {
	return publisher_mapper.MapPublisherToEngine(messageBody, queue, addCorrelationID)
}

func SetMessagePublisher(p MessagePublisher) {
	publisher = p
}
