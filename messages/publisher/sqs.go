package publisher

import (
	"magecomm/logger"
	"magecomm/services"
)

type SQSPublisher struct{}

func (publisher *SQSPublisher) PublishMessage(queueName string, message string) error {
	sqsConnection := services.NewSQSConnection()
	err := sqsConnection.Connect()
	if err != nil {
		logger.Fatalf("Error connecting to SQS: %v", err)
		return err
	}
	sqsClient := sqsConnection.Client

	services.PublishSqsMessage(sqsClient, queueName, message)
	return nil
}
