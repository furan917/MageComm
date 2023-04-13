package publisher

import (
	"fmt"
	"magecomm/messages/listener"
	"magecomm/services"
)

type SQSPublisher struct{}

func (publisher *SQSPublisher) PublishMessage(message string, queueName string, addCorrelationID string) (string, error) {
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

func (publisher *SQSPublisher) GetOutputReturn(correlationID string, queueName string) (string, error) {
	correlationListenerClass := &listener.SqsListener{}
	output, err := correlationListenerClass.ListenForOutputByCorrelationID(queueName, correlationID)
	if err != nil {
		return "", err
	}

	return output, nil
}
