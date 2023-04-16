package listener

import (
	"fmt"
	"magecomm/config_manager"
	"magecomm/logger"
	"magecomm/services"
)

func MapListenerToEngine() (Listener, error) {
	engine := config_manager.GetEngine()
	logger.Debugf("Mapping message to engine:", engine)

	var listenerClass Listener

	switch engine {
	case services.EngineSQS:
		listenerClass = &SqsListener{}
	case services.EngineRabbitMQ:
		listenerClass = &RmqListener{
			ChannelPool: services.RmqChannelPool,
		}
	default:
		return nil, fmt.Errorf("Invalid engine specified for listener: '%s'. Supported engines are: '%s', '%s'.\n", engine, services.EngineSQS, services.EngineRabbitMQ)
	}

	return listenerClass, nil
}
