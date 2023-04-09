package services

import (
	"fmt"
	"github.com/streadway/amqp"
	"magecomm/config_manager"
	"magecomm/logger"
)

const (
	EngineRabbitMQ      = "rmq"
	ConfigRabbitMQTLS   = "rmq_tls"
	ConfigRabbitMQUser  = "rmq_user"
	ConfigRabbitMQPass  = "rmq_pass"
	ConfigRabbitMQHost  = "rmq_host"
	ConfigRabbitMQPort  = "rmq_port"
	ConfigRabbitMQVhost = "rmq_vhost"
)

type RabbitMQConnection struct {
	Connection *amqp.Connection
}

func NewRabbitMQConnection() *RabbitMQConnection {
	return &RabbitMQConnection{}
}

func GetRmqQueueNameWithConfigPrefix(queueName string) string {
	return fmt.Sprintf("%s-%s", config_manager.GetValue(config_manager.CommandConfigEnvironment), queueName)
}

func CreateRmqQueue(channel *amqp.Channel, queueName string) (string, error) {
	//The prefixed name is only used for actual communication, for internal use we use the original name
	queueNameWithConfigPrefix := GetRmqQueueNameWithConfigPrefix(queueName)
	//declare quorum queue
	_, err := channel.QueueDeclare(
		queueNameWithConfigPrefix,
		true,
		false,
		false,
		true,
		amqp.Table{
			"x-queue-type": "quorum",
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to declare quorum queue %q: %v", queueNameWithConfigPrefix, err)
	}

	return queueNameWithConfigPrefix, nil
}

func PublishRmqMessage(channel *amqp.Channel, queueName string, message []byte, messageHeaders amqp.Table) {
	_, err := CreateRmqQueue(channel, queueName)
	if err != nil {
		logger.Fatalf("Failed to create queue: %v", err)
	}
	err = channel.Publish(
		"",
		GetRmqQueueNameWithConfigPrefix(queueName),
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         message,
			Headers:      messageHeaders,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		logger.Warnf("Failed to requeue message: %v", err)
	}
}

func getRabbitMQURL() string {
	useTLS := config_manager.GetValue(ConfigRabbitMQTLS)
	user := config_manager.GetValue(ConfigRabbitMQUser)
	pass := config_manager.GetValue(ConfigRabbitMQPass)
	host := config_manager.GetValue(ConfigRabbitMQHost)
	port := config_manager.GetValue(ConfigRabbitMQPort)
	vhost := config_manager.GetValue(ConfigRabbitMQVhost)

	if user == "" || pass == "" {
		logger.Fatalf("One or more required RabbitMQ environment variables (RMQ_USER, RMQ_PASS) are not set")
	}

	protocol := "amqp"
	switch useTLS {
	case "true", "TRUE", "1":
		protocol = "amqps"
	}

	if vhost[0] != '/' {
		vhost = "/" + vhost
	}

	return fmt.Sprintf("%s://%s:%s@%s:%s%s", protocol, user, pass, host, port, vhost)
}

func CreateRmqChannel() (*RabbitMQConnection, *amqp.Channel, error) {
	rmqConn := NewRabbitMQConnection()
	channel, err := rmqConn.CreateConnectedChannel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create a connected channel: %v", err)
	}

	return rmqConn, channel, nil
}

func (conn *RabbitMQConnection) Connect(rmqURL string) error {
	if conn.Connection != nil {
		return nil
	}

	if rmqURL == "" {
		rmqURL = getRabbitMQURL()
	}

	rmqConn, err := amqp.Dial(rmqURL)
	if err != nil {
		return err
	}
	conn.Connection = rmqConn
	return nil
}

func (conn *RabbitMQConnection) Disconnect() error {
	return conn.Connection.Close()
}

func (conn *RabbitMQConnection) Refresh() error {
	err := conn.Disconnect()
	if err != nil {
		return fmt.Errorf("failed to disconnect RabbitMQ connection: %v", err)
	}

	err = conn.Connect("")
	if err != nil {
		return fmt.Errorf("failed to reconnect to RabbitMQ: %v", err)
	}

	return nil
}

func (conn *RabbitMQConnection) CreateChannel() (*amqp.Channel, error) {
	channel, err := conn.Connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}
	return channel, nil
}

func (conn *RabbitMQConnection) CreateConnectedChannel() (*amqp.Channel, error) {
	err := conn.Connect("")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.CreateChannel()
	if err != nil {
		return nil, fmt.Errorf("failed to create a channel: %v", err)
	}

	return channel, nil
}