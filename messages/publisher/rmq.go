package publisher

import (
	"fmt"
	"github.com/streadway/amqp"
	"magecomm/logger"
	"magecomm/messages/listener"
	"magecomm/services"
)

type RmqPublisher struct{}

func (publisher *RmqPublisher) PublishMessage(message string, queueName string, addCorrelationID string) (string, error) {
	rmqConn, channel, err := services.CreateRmqChannel()
	if err != nil {
		return "", fmt.Errorf("failed to create RabbitMQ channel: %v", err)
	}
	defer func() {
		err := rmqConn.Disconnect()
		if err != nil {
			logger.Warnf("Failed to disconnect from RabbitMQ: %v", err)
		}
	}()
	defer func() {
		err := channel.Close()
		if err != nil {
			logger.Warnf("Failed to close channel: %s", err)
		}
	}()

	correlationID, err := services.PublishRmqMessage(channel, queueName, []byte(message), amqp.Table{}, addCorrelationID)
	return correlationID, nil
}

func (publisher *RmqPublisher) GetOutputReturn(correlationID string, queueName string) (string, error) {
	correlationListenerClass := &listener.RmqListener{
		ChannelPool: services.RmqChannelPool,
	}
	output, err := correlationListenerClass.ListenForOutputByCorrelationID(queueName, correlationID)
	if err != nil {
		return "", err
	}

	return output, nil
}
