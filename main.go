package main

import (
	"github.com/spf13/cobra"
	"magecomm/cmd"
	"magecomm/config_manager"
	"magecomm/logger"
	"magecomm/notifictions"
	"magecomm/services"
)

var RootCmd = &cobra.Command{
	Use:   "magecomm",
	Short: "MageComm CLI is a command line tool for managing Magento applications",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		overrideFilePath, _ := cmd.Flags().GetString("config")
		config_manager.Configure(overrideFilePath)
		initializeModuleWhichRequireConfig()

		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			logger.EnableDebugMode()
		}
	},
}

func initializeModuleWhichRequireConfig() {
	notifictions.Initialize()
	services.InitializeRMQ()
	services.InitializeSQS()
}

func main() {
	RootCmd.AddCommand(cmd.ListenCmd)
	RootCmd.AddCommand(cmd.MagerunCmd)
	RootCmd.AddCommand(cmd.DeployCmd)
	RootCmd.AddCommand(cmd.CatCmd)
	RootCmd.AddCommand(cmd.CatDeployCmd)

	RootCmd.PersistentFlags().String("config", "", "Path to config file")
	RootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")

	err := RootCmd.Execute()
	if err != nil {
		logger.Fatalf("Failed to execute command: %s", err)
		return
	}
}
