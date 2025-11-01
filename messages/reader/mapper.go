package reader

import (
	"fmt"
	"magecomm/config_manager"
	"magecomm/services"
)

var Reader MessageReader

func MapReaderToEngine() (MessageReader, error) {
	engine := config_manager.GetEngine()

	switch engine {
	case services.EngineSQS:
		return &SqsReader{}, nil
	case services.EngineRabbitMQ:
		return &RmqReader{}, nil
	default:
		return nil, fmt.Errorf("unknown engine: %s", engine)
	}
}