package publisher

import (
	"fmt"
	"magecomm/config_manager"
	"magecomm/logger"
	"magecomm/services"
)

func MapPublisherToEngine() (MessagePublisher, error) {
	engine := config_manager.GetEngine()
	logger.Debugf("Mapping message to engine: %v", engine)
	var publisherClass MessagePublisher

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
