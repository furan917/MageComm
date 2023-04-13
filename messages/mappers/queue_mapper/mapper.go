package queue_mapper

import (
	"fmt"
	"magecomm/config_manager"
)

func MapQueueToOutputQueue(queueName string) string {
	return fmt.Sprintf("%s_%s", queueName, "output")
}

func MapQueueToEngineOutputQueue(queueName string) string {
	queueName = MapQueueToOutputQueue(queueName)
	return MapQueueToEngineQueue(queueName)
}

func MapQueueToEngineQueue(queueName string) string {
	return fmt.Sprintf("%s-%s", config_manager.GetValue(config_manager.CommandConfigEnvironment), queueName)
}