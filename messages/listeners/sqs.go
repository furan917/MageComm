package listeners

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"magecomm/logger"
	"magecomm/messages/handler"
	"magecomm/services"
	"strconv"
	"time"
)

func processSqsMessages(sqsClient *sqs.SQS, queueName string) {
	queueURL, err := services.CreateSQSQueueIfNotExists(sqsClient, queueName)
	if err != nil {
		queueNameWithConfigPrefix := services.GetSqsQueueNameWithConfigPrefix(queueName)
		logger.Fatalf("Error building SQS queue URL for queue %s: %v\n", queueNameWithConfigPrefix, err)
		return
	}

	for {
		result, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: aws.Int64(1),
			VisibilityTimeout:   aws.Int64(60),
			WaitTimeSeconds:     aws.Int64(0),
			AttributeNames:      aws.StringSlice([]string{"All"}),
		})

		if err != nil {
			logger.Warnf("Error receiving message from SQS:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if len(result.Messages) == 0 {
			logger.Warnf("No messages available. Waiting...")
			time.Sleep(5 * time.Second)
			continue
		}

		message := result.Messages[0]
		receiveCount, _ := strconv.Atoi(*message.Attributes["ApproximateReceiveCount"])
		messageBody := *message.Body
		logger.Debugf("Message received from", queueURL, ":", *message.Body)
		if err := handler.HandleReceivedMessage(queueName, messageBody); err != nil {
			logger.Warnf("Error handling message, could not process command:", messageBody,
				" retry attempt:", receiveCount, "of 5",
				" error:", err)
			if receiveCount >= handler.MessageRetryLimit {
				continue
			}
		}

		_, err = sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      aws.String(queueURL),
			ReceiptHandle: message.ReceiptHandle,
		})

		if err != nil {
			logger.Warnf("Error deleting message from SQS:", err)
		}
	}
}

func listenToQueue(queueName string) {

	sqsConnection := services.NewSQSConnection()
	err := sqsConnection.Connect()
	if err != nil {
		logger.Fatalf("Error connecting to SQS: %v", err)
	}
	sqsClient := sqsConnection.Client
	processSqsMessages(sqsClient, queueName)
}

func StartListeningSqs(queueNames []string) {
	for _, queueName := range queueNames {
		go listenToQueue(queueName)
	}

	// Wait indefinitely to prevent the program from exiting.
	select {}
}
