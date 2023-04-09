package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var DeployCmd = &cobra.Command{
	Use:   "deploy [filepath]",
	Short: "Deploy a gzipped file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		// Implement your deployment logic here
		fmt.Printf("Deploying file: %s\n", filepath)
	},
}
