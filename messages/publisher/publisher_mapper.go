package publisher

import (
	"fmt"
	"magecomm/config_manager"
	"magecomm/loading"
	"magecomm/logger"
	"magecomm/services"
	"sync"
)

func MapPublisherToEngine() (Publisher, error) {
	engine := config_manager.GetEngine()
	logger.Debugf("Mapping message to engine:", engine)
	var publisherClass Publisher

	switch engine {
	case services.EngineSQS:
		publisherClass = &SQSPublisher{}
	case services.EngineRabbitMQ:
		publisherClass = &RmqPublisher{}
	default:
		return nil, fmt.Errorf("Invalid engine specified for publisher: '%s'. Supported engines are: '%s','%s'.\n", engine, services.EngineSQS, services.EngineRabbitMQ)
	}

	return publisherClass, nil
}

func HandleCorrelationID(publisherClass Publisher, correlationID string, queueName string) (string, error) {
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
