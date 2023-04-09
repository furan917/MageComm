package publisher

type Publisher interface {
	PublishMessage(queueName string, Message string) error
}
