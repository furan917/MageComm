package publisher

import (
	"magecomm/logger"
	"magecomm/services"
)

func PublishSqsMessage(queueName string, messageBody string) {
	sqsConnection := services.NewSQSConnection()
	err := sqsConnection.Connect()
	if err != nil {
		logger.Fatalf("Error connecting to SQS: %v", err)
	}
	sqsClient := sqsConnection.Client

	services.PublishSqsMessage(sqsClient, queueName, messageBody)
}
