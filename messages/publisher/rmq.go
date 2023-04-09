package publisher

import (
	"github.com/streadway/amqp"
	"magecomm/logger"
	"magecomm/services"
)

func PublishRmqMessage(queueName string, message string) {
	rmqConn, channel, err := services.CreateRmqChannel()
	if err != nil {
		logger.Fatalf("Failed to create RabbitMQ channel: %v", err)
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

	services.PublishRmqMessage(channel, queueName, []byte(message), amqp.Table{})
}
