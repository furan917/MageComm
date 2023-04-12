package main

import (
	"github.com/spf13/cobra"
	"magecomm/cmd"
	"magecomm/config_manager"
	"magecomm/logger"
)

var RootCmd = &cobra.Command{
	Use:   "magecomm",
	Short: "MageComm CLI is a command line tool for managing Magento applications",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			logger.EnableDebugMode()
		}
	},
}

func main() {
	config_manager.Configure()

	RootCmd.AddCommand(cmd.ListenCmd)
	RootCmd.AddCommand(cmd.MagerunCmd)
	RootCmd.AddCommand(cmd.DeployCmd)
	RootCmd.AddCommand(cmd.CatCmd)
	RootCmd.AddCommand(cmd.CatDeployCmd)

	RootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")

	err := RootCmd.Execute()
	if err != nil {
		logger.Fatalf("Failed to execute command: %s", err)
		return
	}
}
