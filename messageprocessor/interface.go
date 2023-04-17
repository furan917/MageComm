package messageprocessor

type MessageProcessor interface {
	ProcessMessage(message string, CorrelationId string) error
}
