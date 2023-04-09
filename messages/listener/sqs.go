package listener

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"magecomm/logger"
	"magecomm/messages/handler"
	"magecomm/services"
	"magecomm/system_limits"
	"strconv"
	"sync"
	"time"
)

type SqsListener struct {
	stopChan  chan struct{}
	waitGroup sync.WaitGroup
}

func (listener *SqsListener) shouldExecutionBeDelayed() {
	for system_limits.CheckIfOutsideOperationalLimits() {
		system_limits.SystemLimitCheckSleep()
	}
}

func (listener *SqsListener) processSqsMessage(message *sqs.Message, sqsClient *sqs.SQS, queueName string, queueURL string) {
	receiveCount, err := strconv.Atoi(*message.Attributes["ApproximateReceiveCount"])
	if err != nil {
		logger.Warnf("Error parsing ApproximateReceiveCount attribute: %v", err)
	}

	messageBody := *message.Body
	logger.Debugf("Message received from", queueName)

	listener.shouldExecutionBeDelayed()
	if err := handler.HandleReceivedMessage(queueName, messageBody); err != nil {
		logger.Warnf("Error handling message, could not process command:", messageBody,
			" retry attempt:", receiveCount, "of 5",
			" error:", err)
		if receiveCount <= handler.MessageRetryLimit {
			return
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

// This loop is indefinite and contains an anonymous function to ensure timeouts are handled correctly
func (listener *SqsListener) loopThroughMessages(sqsClient *sqs.SQS, queueName string, queueURL string) {
	for {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			result, err := sqsClient.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(queueURL),
				MaxNumberOfMessages: aws.Int64(1),
				VisibilityTimeout:   aws.Int64(60),
				WaitTimeSeconds:     aws.Int64(0),
				AttributeNames:      aws.StringSlice([]string{"All"}),
			})

			if err != nil {
				logger.Warnf("Error receiving message from SQS:", err)
				time.Sleep(5 * time.Second)
				return
			}

			if len(result.Messages) == 0 {
				logger.Warnf("No messages available. Waiting...")
				time.Sleep(5 * time.Second)
				return
			}

			for _, message := range result.Messages {
				listener.processSqsMessage(message, sqsClient, queueName, queueURL)
			}
		}()
	}
}

func (listener *SqsListener) listenToQueue(queueName string) {
	listener.waitGroup.Add(1)
	defer listener.waitGroup.Done()

	sqsConnection, err := services.GetSQSConnection()
	if err != nil {
		logger.Fatalf("Unable to get SQS connection from pool %v", err)
	}

	defer services.ReleaseSQSConnection(sqsConnection)

	err = sqsConnection.Connect()
	if err != nil {
		logger.Fatalf("Error connecting to SQS: %v", err)
	}
	sqsClient := sqsConnection.Client

	queueURL, err := services.CreateSQSQueueIfNotExists(sqsClient, queueName)
	if err != nil {
		queueNameWithConfigPrefix := services.GetSqsQueueNameWithConfigPrefix(queueName)
		logger.Fatalf("Error building SQS queue URL for queue %s: %v\n", queueNameWithConfigPrefix, err)
		return
	}

	for {
		select {
		case <-listener.stopChan:
			return
		default:
			listener.loopThroughMessages(sqsClient, queueName, queueURL)
		}
	}
}

func (listener *SqsListener) ListenToService(queueNames []string) {
	listener.stopChan = make(chan struct{})
	for _, queueName := range queueNames {
		go listener.listenToQueue(queueName)
	}

	// Wait indefinitely to prevent the program from exiting.
	select {}
}

func (listener *SqsListener) Close() {
	close(listener.stopChan)
	listener.waitGroup.Wait()
}