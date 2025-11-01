package reader

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"magecomm/logger"
	"magecomm/messages/queues"
	"magecomm/services"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SqsReader struct{}

func (r *SqsReader) DrainOutputQueue(queueName string) (int, error) {
	outputQueueName := queues.MapQueueToOutputQueue(queueName)

	sqsConnection, err := services.GetSQSConnection()
	if err != nil {
		return 0, fmt.Errorf("unable to get SQS connection: %w", err)
	}
	defer services.ReleaseSQSConnection(sqsConnection)

	if err := sqsConnection.Connect(); err != nil {
		return 0, fmt.Errorf("error connecting to SQS: %w", err)
	}

	sqsClient := sqsConnection.Client
	queueURL, err := services.CreateSQSQueueIfNotExists(sqsClient, outputQueueName)
	if err != nil {
		return 0, fmt.Errorf("error getting queue URL: %w", err)
	}

	count := 0
	for {
		result, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:              aws.String(queueURL),
			MaxNumberOfMessages:   aws.Int64(10),
			WaitTimeSeconds:       aws.Int64(1),
			AttributeNames:        []*string{aws.String("SentTimestamp")},
			MessageAttributeNames: []*string{aws.String("All")},
		})

		if err != nil {
			return count, fmt.Errorf("error receiving messages: %w", err)
		}

		if len(result.Messages) == 0 {
			break
		}

		for _, msg := range result.Messages {
			count++
			r.displayMessage(count, msg)

			_, err := sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				logger.Warnf("Failed to delete message: %v", err)
			}
		}
	}

	return count, nil
}

func (r *SqsReader) displayMessage(index int, msg *sqs.Message) {
	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Printf("Output %d\n", index)

	if msg.Attributes["SentTimestamp"] != nil {
		timestamp := *msg.Attributes["SentTimestamp"]
		if ms, err := strconv.ParseInt(timestamp, 10, 64); err == nil {
			t := time.Unix(ms/1000, 0)
			fmt.Printf("Received: %s\n", t.Format("2006-01-02 15:04:05"))
		}
	}

	fmt.Printf("%s\n", strings.Repeat("=", 60))
	fmt.Println(*msg.Body)
}
