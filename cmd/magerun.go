package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"magecomm/config_manager"
	"magecomm/loading"
	"magecomm/logger"
	"magecomm/messages/listener"
	"magecomm/messages/publisher"
	"strings"
	"sync"
)

const MageRunQueue = "magerun"

func handleMageRunCmdMessage(args []string) error {
	messageBody := strings.Join(args, " ")
	publisherClass := publisher.Publisher
	correlationID, err := publisherClass.Publish(messageBody, MageRunQueue, uuid.New().String())
	if err != nil {
		return fmt.Errorf("failed to publish message: %s", err)
	}

	if correlationID == "" {
		fmt.Println("Command executed, but no output could be returned")
		return nil
	}

	output, err := handleCorrelationID(correlationID, MageRunQueue)
	if err != nil {
		return fmt.Errorf("failed to get output: %s", err)
	}

	if output != "" {
		fmt.Println(output)
	} else {
		fmt.Println("Command executed, but no output was returned")
	}

	return nil
}

var MagerunCmd = &cobra.Command{
	Use:                "magerun",
	Short:              "A wrapper for the magerun command with restricted command usage",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("no command provided")
		}

		command := args[0]
		if !config_manager.IsMageRunCommandAllowed(command) {
			return fmt.Errorf("the command '%s' is not allowed", command)
		}

		err := handleMageRunCmdMessage(args)
		if err != nil {
			return err
		}
		return nil
	},
}

func handleCorrelationID(correlationID string, queueName string) (string, error) {
	logger.Debugf("Correlation ID:", correlationID, "returned from publisher. Listening for output.")

	var wg sync.WaitGroup
	stopLoading := make(chan bool)
	wg.Add(1)
	go func() {
		defer wg.Done()
		loading.Indicator(stopLoading)
	}()

	listenerClass := listener.Listener
	output, err := listenerClass.ListenForOutputByCorrelationID(queueName, correlationID)
	if err != nil {
		return "", err
	}

	// Stop the loading indicator
	stopLoading <- true
	wg.Wait()

	return output, err
}
