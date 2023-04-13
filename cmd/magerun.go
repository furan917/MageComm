package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"magecomm/config_manager"
	"strings"
)

const MageRunQueue = "magerun"

func handleMageRunCmdMessage(args []string) error {
	messageBody := strings.Join(args, " ")
	output, err := publisher.Publish(messageBody, MageRunQueue, uuid.New().String())
	if err != nil {
		return fmt.Errorf("failed to publish message: %s", err)
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
