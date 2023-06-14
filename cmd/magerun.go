package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"magecomm/config_manager"
	"magecomm/messages/listener"
	"magecomm/messages/publisher"
	"strings"
)

const MageRunQueue = "magerun"

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

		if config_manager.IsRestrictedCommandArgsIncluded(command, args[1:]) {
			return fmt.Errorf("the command '%s' is not allowed with the following arguments: %s", command, strings.Join(args[1:], " "))
		}

		if isAllRequiredArgsIncluded, missingRequiredArgs := config_manager.IsRequiredCommandArgsIncluded(command, args[1:]); !isAllRequiredArgsIncluded {
			prompt := fmt.Sprintf("The command '%s' is missing required arguments: %s. Do you want to run this command and include them?", command, strings.Join(missingRequiredArgs, " "))
			confirmed, err := PromptUserForConfirmation(prompt)
			if err != nil {
				return fmt.Errorf("error while reading user input: %v", err)
			}

			if confirmed {
				args = append(args, missingRequiredArgs...)
			} else {
				return fmt.Errorf("exiting: the command '%s' cannot be executed without the required arguments: %s", command, strings.Join(missingRequiredArgs, " "))
			}
		}

		err := handleMageRunCmdMessage(args)
		if err != nil {
			return err
		}
		return nil
	},
}

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

	output, err := listener.HandleOutputByCorrelationID(correlationID, MageRunQueue)
	if err != nil {
		return fmt.Errorf("failed to get output: %s", err)
	}

	if output != "" {
		fmt.Println(output)
	}

	return nil
}
