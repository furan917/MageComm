package listener

import (
	"fmt"
	"magecomm/config_manager"
	"magecomm/config_manager/loading"
	"magecomm/logger"
	"strconv"
	"sync"
	"time"
)

func HandleOutputByCorrelationID(correlationID string, queueName string) (string, error) {
	logger.Debugf("Listening for output with Correlation ID:", correlationID)

	var wg sync.WaitGroup
	stopLoading := make(chan bool)
	wg.Add(1)
	go func() {
		defer wg.Done()
		loading.Indicator(stopLoading)
	}()

	listenerClass := Listener
	outputCh := make(chan string)
	errCh := make(chan error)

	var output string
	go func() {
		output, err := listenerClass.ListenForOutputByCorrelationID(queueName, correlationID)
		if err != nil {
			errCh <- err
			return
		}
		outputCh <- output
	}()

	timeout, err := strconv.Atoi(config_manager.GetValue(config_manager.CommandConfigPublisherOutputTimeout))
	if err != nil {
		timeout = 60
	}
	timeoutDuration := time.Duration(timeout) * time.Second

	select {
	case output = <-outputCh:
	case err = <-errCh:
	case <-time.After(timeoutDuration):
		err = fmt.Errorf("waiting for command timed out after %v", timeoutDuration)
	}

	stopLoading <- true
	wg.Wait()

	return output, err
}
