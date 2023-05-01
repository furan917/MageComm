package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"magecomm/logger"
	"magecomm/messages/listener"
	"magecomm/messages/publisher"
)

const DeployQueue = "deploy"

var DeployCmd = &cobra.Command{
	Use:   "deploy [filename]",
	Short: "Deploy a latest version of site",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fileName := args[0]
		fmt.Printf("Deploying file: %s\n", fileName)
		err := handleDeployCmdMessage(fileName)
		if err != nil {
			logger.Fatal(err)
		}
	},
}

func handleDeployCmdMessage(fileName string) error {
	messageBody := fileName
	publisherClass := publisher.Publisher
	correlationID, err := publisherClass.Publish(messageBody, DeployQueue, uuid.New().String())
	if err != nil {
		return fmt.Errorf("failed to publish deploy message: %s", err)
	}

	if correlationID == "" {
		fmt.Println("deploy executed, but no output could be returned")
		return nil
	}

	output, err := listener.HandleOutputByCorrelationID(correlationID, DeployQueue)
	if err != nil {
		return fmt.Errorf("failed to get output: %s", err)
	}

	if output != "" {
		fmt.Println(output)
	} else {
		fmt.Println("deploy executed, but no output was returned")
	}

	return nil
}
