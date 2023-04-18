package publisher

import (
	"fmt"
	"magecomm/services"
)

type SQSPublisher struct{}

func (publisher *SQSPublisher) Publish(message string, queueName string, addCorrelationID string) (string, error) {
	sqsConnection := services.NewSQSConnection()
	err := sqsConnection.Connect()
	if err != nil {
		return "", fmt.Errorf("error connecting to SQS: %v", err)
	}
	sqsClient := sqsConnection.Client

	correlationID, err := services.PublishSqsMessage(sqsClient, queueName, message, addCorrelationID)
	if err != nil {
		return "", fmt.Errorf("error publishing message to SQS: %v", err)
	}
	return correlationID, nil
}
