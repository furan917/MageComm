package messages

import (
	"log"
	"magecomm/config_manager"
	"magecomm/logger"
	"magecomm/messages/listener"
	"magecomm/messages/publisher"
	"magecomm/services"
	"os"
	"os/signal"
	"syscall"
)

func MapListenerToEngine(queueNames []string) {
	engine := getEngine()
	logger.Debugf("Listening to queues:", queueNames, "using engine:", engine)

	var listenerClass listener.Listener

	switch engine {
	case services.EngineSQS:
		listenerClass = &listener.SqsListener{}
	case services.EngineRabbitMQ:
		listenerClass = &listener.RmqListener{
			ChannelPool: services.RmqChannelPool,
		}
	default:
		logger.Fatalf("Invalid engine specified for listener: '%s'. Supported engines are: '%s', '%s'.\n", engine, services.EngineSQS, services.EngineRabbitMQ)
		return
	}

	// Create a channel to handle program termination or interruption signals so we can kill any connections if needed
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go listenerClass.ListenToService(queueNames)
	<-sigChan
	listenerClass.Close()
}

func MapPublisherToEngine(queueName string, messageBody string) {
	engine := getEngine()
	logger.Debugf("publishing message:", messageBody, " to queue: ", queueName, "on engine:", engine)

	var publishClass publisher.Publisher

	switch engine {
	case services.EngineSQS:
		publishClass = &publisher.SQSPublisher{}
	case services.EngineRabbitMQ:
		publishClass = &publisher.RmqPublisher{}
	default:
		logger.Fatalf("Invalid engine specified for publisher: '%s'. Supported engines are: '%s','%s'.\n", engine, services.EngineSQS, services.EngineRabbitMQ)
		return
	}

	err := publishClass.PublishMessage(messageBody, queueName)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}
}

func getEngine() string {
	return config_manager.GetValue(config_manager.CommandConfigListenerEngine)
}
