package publisher_mapper

import (
	"fmt"
	"magecomm/config_manager"
	"magecomm/loading"
	"magecomm/logger"
	"magecomm/messages/publisher"
	"magecomm/services"
	"sync"
)

func MapPublisherToEngine(messageBody string, queueName string, addCorrelationID string) (string, error) {
	engine := getEngine()
	logger.Debugf("publishing message:", messageBody, " to queue: ", queueName, "on engine:", engine)

	var publisherClass publisher.Publisher

	switch engine {
	case services.EngineSQS:
		publisherClass = &publisher.SQSPublisher{}
	case services.EngineRabbitMQ:
		publisherClass = &publisher.RmqPublisher{}
	default:
		return "", fmt.Errorf("Invalid engine specified for publisher: '%s'. Supported engines are: '%s','%s'.\n", engine, services.EngineSQS, services.EngineRabbitMQ)
	}

	correlationID, err := publisherClass.PublishMessage(messageBody, queueName, addCorrelationID)
	if err != nil {
		return "", fmt.Errorf("failed to publish message: %v", err)
	}
	if correlationID != "" {
		return handleCorrelationID(publisherClass, correlationID, queueName)
	}

	return "", nil
}

// todo this causes a cyclic import =: listener > handler > magerun > publisher ?? listener again due to the handle correlation function call above
// this should be moved outside of package and only published cmd messages, not output messages should use this: otherwise we create an infinite loop
func handleCorrelationID(publisherClass publisher.Publisher, correlationID string, queueName string) (string, error) {
	logger.Debugf("Correlation ID:", correlationID, "returned from publisher. Listening for output.")

	var wg sync.WaitGroup
	stopLoading := make(chan bool)
	wg.Add(1)
	go func() {
		defer wg.Done()
		loading.Indicator(stopLoading)
	}()

	output, err := publisherClass.GetOutputReturn(correlationID, queueName)

	// Stop the loading indicator
	stopLoading <- true
	wg.Wait()

	return output, err
}

func getEngine() string {
	return config_manager.GetValue(config_manager.CommandConfigListenerEngine)
}
