package listeners

import (
	"github.com/streadway/amqp"
	"magecomm/logger"
	"magecomm/messages/handler"
	"magecomm/services"
)

func processRmqMessages(channel *amqp.Channel, queueName string) {
	queueNameWithConfigPrefix, err := services.CreateRmqQueue(channel, queueName)
	if err != nil {
		return
	}
	msgs, err := channel.Consume(
		queueNameWithConfigPrefix,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Fatalf("%s: %s", "Failed to register a consumer", err)
	}

	for msg := range msgs {
		logger.Warnf("Received a message: %s", msg.Body)
		retryCount, ok := msg.Headers["RetryCount"]
		if !ok {
			retryCount = 0
		}

		if err := handler.HandleReceivedMessage(queueName, string(msg.Body)); err != nil {
			logger.Warnf("Failed to process message: %v", err)
			if retryCount.(int) < handler.MessageRetryLimit {
				msg.Headers["RetryCount"] = retryCount.(int) + 1
				services.PublishRmqMessage(channel, queueName, msg.Body, msg.Headers)
			} else {
				logger.Warnf("Retry count exceeded. Discarding the message.")
			}
		}
	}
}

func listenToRabbitMQQueue(queueName string) {
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

	processRmqMessages(channel, queueName)
}

func StartListeningRmq(queueNames []string) {
	for _, queueName := range queueNames {
		go listenToRabbitMQQueue(queueName)
	}

	// Wait indefinitely to prevent the program from exiting.
	select {}
}
