package messages

import (
	"fmt"
	"magecomm/config_manager"
	"magecomm/logger"
	"magecomm/messages/listeners"
	"magecomm/messages/publisher"
	"magecomm/services"
)

func MapListenerToEngine(queueNames []string) {
	engine := getEngine()
	logger.Debugf("Listening to queues:", queueNames, "using engine:", engine)

	switch engine {
	case services.EngineSQS:
		listeners.StartListeningSqs(queueNames)
	case services.EngineRabbitMQ:
		listeners.StartListeningRmq(queueNames)
	default:
		logger.Fatalf("Invalid engine specified for listener: '%s'. Supported engines are: '%s', '%s'.\n", engine, services.EngineSQS, services.EngineRabbitMQ)
		return
	}
}

func MapPublisherToEngine(queueName string, messageBody string) {
	engine := getEngine()
	fmt.Println("publishing message:", messageBody, " to queue: ", queueName, "on engine:", engine)

	switch engine {
	case services.EngineSQS:
		publisher.PublishSqsMessage(queueName, messageBody)
	case services.EngineRabbitMQ:
		publisher.PublishRmqMessage(queueName, messageBody)
	default:
		logger.Fatalf("Invalid engine specified for publisher: '%s'. Supported engines are: '%s','%s'.\n", engine, services.EngineSQS, services.EngineRabbitMQ)
		return
	}
}

func getEngine() string {
	return config_manager.GetValue("listener_engine")
}
