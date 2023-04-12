package listener

import (
    "errors"
    "github.com/streadway/amqp"
    "magecomm/logger"
    "magecomm/messages/handler"
    "magecomm/services"
    "magecomm/system_limits"
    "sync"
)

type RmqListener struct {
	ChannelPool *services.RabbitMQChannelPool
	done        chan struct{}
	wg          sync.WaitGroup
}

func (listener *RmqListener) shouldExecutionBeDelayed() error {
	totalDeferTime := 0
	for system_limits.CheckIfOutsideOperationalLimits() {
		system_limits.SystemLimitCheckSleep()
		totalDeferTime += int(system_limits.WaitTimeBetweenChecks)

		if totalDeferTime > int(system_limits.MaxDeferralTime) {
			return errors.New("max deferral time exceeded")
		}
	}

	return nil
}

func (listener *RmqListener) processRmqMessage(message amqp.Delivery, channel *amqp.Channel, queueName string) {
	logger.Debugf("Message received from", queueName)
	retryCount, ok := message.Headers["RetryCount"]
	if !ok {
		retryCount = 0
	}

	err := listener.shouldExecutionBeDelayed()
	if err != nil {
		logger.Warnf("Message deferral time exceeded. Dropping hold on the message.")
		message.Headers["RetryCount"] = retryCount.(int) + 1
		services.PublishRmqMessage(channel, queueName, message.Body, message.Headers)
		return
	}
	if err := handler.HandleReceivedMessage(queueName, string(message.Body)); err != nil {
		logger.Warnf("Failed to process message: %v", err)
		if retryCount.(int) < handler.MessageRetryLimit {
			message.Headers["RetryCount"] = retryCount.(int) + 1
			services.PublishRmqMessage(channel, queueName, message.Body, message.Headers)
		} else {
			logger.Warnf("Retry count exceeded. Discarding the message.")
		}
	}
}

func (listener *RmqListener) listenToQueue(queueName string) {
	defer listener.wg.Done()

	channel, err := listener.ChannelPool.Get()
	if err != nil {
		logger.Warnf("Error getting channel from pool: %v", err)
		return
	}
	defer listener.ChannelPool.Put(channel)

	queueNameWithConfigPrefix, err := services.CreateRmqQueue(channel, queueName)
	if err != nil {
		return
	}
	msgs, err := channel.Consume(
		queueNameWithConfigPrefix,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Fatalf("%s: %s", "Failed to register a consumer", err)
	}

	for {
		select {
		case message, ok := <-msgs:
			if !ok {
				return
			}
			listener.processRmqMessage(message, channel, queueName)
		case <-listener.done:
			return
		}
	}
}

func (listener *RmqListener) ListenToService(queueNames []string) {
	listener.done = make(chan struct{})

	for _, queueName := range queueNames {
		listener.wg.Add(1)
		go listener.listenToQueue(queueName)
	}

	listener.wg.Wait()
}

func (listener *RmqListener) Close() {
	close(listener.done)
}
