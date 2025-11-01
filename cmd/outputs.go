package cmd

import (
	"fmt"
	"magecomm/messages/reader"

	"github.com/spf13/cobra"
)

var OutputsCmd = &cobra.Command{
	Use:   "outputs",
	Short: "Drain all pending command outputs from the queue",
	RunE: func(cmd *cobra.Command, args []string) error {
		readerInstance, err := reader.MapReaderToEngine()
		if err != nil {
			return err
		}

		count, err := readerInstance.DrainOutputQueue("magerun")
		if err != nil {
			return err
		}

		if count == 0 {
			fmt.Println("No outputs available")
		} else {
			fmt.Println("\nOutputs finished")
		}

		return nil
	},
}
