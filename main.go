package main

import (
	"github.com/spf13/cobra"
	"magecomm/cmd"
	"magecomm/config_manager"
	"magecomm/logger"
)

func main() {
	config_manager.Configure()

	var rootCmd = &cobra.Command{
		Use:   "magecomm",
		Short: "MageComm CLI is a command line tool for managing Magento applications",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			debug, _ := cmd.Flags().GetBool("debug")
			if debug {
				logger.EnableDebugMode()
			}
		},
	}
	rootCmd.AddCommand(cmd.ListenCmd)
	rootCmd.AddCommand(cmd.MagerunCmd)
	rootCmd.AddCommand(cmd.DeployCmd)
	rootCmd.AddCommand(cmd.CatCmd)
	rootCmd.AddCommand(cmd.CatDeployCmd)

	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")

	err := rootCmd.Execute()
	if err != nil {
		logger.Fatalf("Failed to execute command: %s", err)
		return
	}
}
