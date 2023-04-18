package handler

type MessageHandler interface {
	ProcessMessage(message string, CorrelationId string) error
}
