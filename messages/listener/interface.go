package listener

import "fmt"

type MessageListener interface {
	ListenToService(queueNames []string)
	ListenForOutputByCorrelationID(queueName string, correlationID string) (string, error)
	Close()
}

var Listener MessageListener = &DefaultMessageListener{}

type DefaultMessageListener struct{}

func (d *DefaultMessageListener) ListenToService(queueNames []string) {
	listener, err := MapListenerToEngine()
	if err != nil {
		fmt.Printf("failed to map listener to engine: %v\n", err)
		return
	}

	listener.ListenToService(queueNames)
}

func (d *DefaultMessageListener) ListenForOutputByCorrelationID(queueName string, correlationID string) (string, error) {
	listener, err := MapListenerToEngine()
	if err != nil {
		return "", fmt.Errorf("failed to map listener to engine: %v", err)
	}

	output, err := listener.ListenForOutputByCorrelationID(queueName, correlationID)
	if err != nil {
		return "", fmt.Errorf("failed to get output by correlation ID: %v", err)
	}

	return output, nil
}

func (d *DefaultMessageListener) Close() {
	listener, err := MapListenerToEngine()
	if err != nil {
		fmt.Printf("failed to map listener to engine: %v\n", err)
		return
	}

	listener.Close()
}

func SetMessageListener(l MessageListener) {
	Listener = l
}
