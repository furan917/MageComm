package cmd

import "magecomm/messages"

type MessagePublisher interface {
	Publish(queue string, messageBody string)
}

var publisher MessagePublisher = &defaultMessagePublisher{}

type defaultMessagePublisher struct{}

func (d *defaultMessagePublisher) Publish(queue string, messageBody string) {
	messages.MapPublisherToEngine(queue, messageBody)
}

func SetMessagePublisher(p MessagePublisher) {
	publisher = p
}
