package reader

type MessageReader interface {
	DrainOutputQueue(queueName string) (int, error)
}