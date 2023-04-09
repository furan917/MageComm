package listener

type Listener interface {
	ListenToService(queueNames []string)
	listenToQueue(queueName string)
	shouldExecutionBeDelayed()
	Close()
}
