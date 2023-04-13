package listener

type Listener interface {
	ListenToService(queueNames []string)
	listenToQueue(queueName string)
	ListenForOutputByCorrelationID(queueName string, correlationID string) (string, error)
	shouldExecutionBeDelayed() error
	Close()
}
