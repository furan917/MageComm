package publisher

type Publisher interface {
	PublishMessage(Message string, queueName string, addCorrelationID string) (string, error)
	GetOutputReturn(correlationID string, queueName string) (string, error)
}
