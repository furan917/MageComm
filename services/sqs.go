package services

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sts"
	"log"
	"magecomm/config_manager"
	"magecomm/logger"
	"magecomm/messages/mappers/queue_mapper"
)

const (
	EngineSQS = "sqs"
)

type SQSConnection struct {
	Client *sqs.SQS
}

func NewSQSConnection() *SQSConnection {
	return &SQSConnection{}
}

func CreateSQSQueueIfNotExists(sqsClient *sqs.SQS, queueName string) (string, error) {
	//The prefixed name is only used for actual communication, for internal use we use the original name
	queueName = queue_mapper.MapQueueToEngineQueue(queueName)
	getQueueUrlInput := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}
	_, err := sqsClient.GetQueueUrl(getQueueUrlInput)
	if err == nil {
		return BuildSQSQueueURL(sqsClient, queueName)
	}

	createQueueInput := &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	}
	createQueueOutput, err := sqsClient.CreateQueue(createQueueInput)
	if err != nil {
		return "", fmt.Errorf("failed to create SQS queue: %v", err)
	}

	return aws.StringValue(createQueueOutput.QueueUrl), nil
}

func BuildSQSQueueURL(sqsClient *sqs.SQS, queueName string) (string, error) {
	awsRegion := aws.StringValue(sqsClient.Config.Region)
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      sqsClient.Config.Region,
		Credentials: sqsClient.Config.Credentials,
	}))
	stsSvc := sts.New(sess)
	result, err := stsSvc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", fmt.Errorf("failed to get AWS account ID: %v", err)
	}
	awsAccountID := aws.StringValue(result.Account)

	return fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/%s", awsRegion, awsAccountID, queueName), nil
}

func PublishSqsMessage(sqsClient *sqs.SQS, queueName string, messageBody string, addCorrelationID string) (string, error) {
	correlationID := ""
	if addCorrelationID != "" {
		correlationID = addCorrelationID
	}

	queueURL, err := CreateSQSQueueIfNotExists(sqsClient, queueName)
	_, err = sqsClient.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(messageBody),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"CorrelationID": {
				DataType:    aws.String("String"),
				StringValue: aws.String(correlationID),
			},
		},
	})

	if err != nil {
		return "", fmt.Errorf("failed to publish message: %v", err)
	}
	logger.Debugf("Message published successfully with correlation ID: %s", correlationID)
	return correlationID, nil

}

// Connect to SQS using IAM role
func (conn *SQSConnection) Connect() error {
	awsRegion := config_manager.GetValue(config_manager.ConfigSQSRegion)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if err != nil {
		log.Printf("Error connecting to SQS using IAM role: %v", err)
		return err
	}

	conn.Client = sqs.New(sess)
	return nil
}

func (conn *SQSConnection) Disconnect() {
	conn.Client = nil
	log.Printf("SQS connection disconnected")
}

func (conn *SQSConnection) Refresh() error {
	conn.Disconnect()
	err := conn.Connect()
	if err != nil {
		log.Printf("Error connecting to SQS during refresh: %v", err)
		return err
	}

	log.Printf("SQS connection successfully refreshed")
	return nil
}
