package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"magecomm/magerun"
	"magecomm/messages"
	"strings"
)

const MageRunQueue = "magerun"

func publishMageRunMessage(args []string) {
	messageBody := strings.Join(args, " ")
	messages.MapPublisherToEngine(MageRunQueue, messageBody)
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
		if !magerun.IsCommandAllowed(command) {
			return fmt.Errorf("the command '%s' is not allowed", command)
		}

		publishMageRunMessage(args)
		return nil
	},
}
