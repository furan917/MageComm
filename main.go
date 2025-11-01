package main

import (
	_ "embed"
	"fmt"
	"magecomm/cmd"
	"magecomm/config_manager"
	"magecomm/logger"
	"magecomm/notifictions"
	"magecomm/services"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed version.txt
var version string

var RootCmd = &cobra.Command{
	Use:   "magecomm",
	Short: "MageComm CLI is a command line tool for managing Magento applications",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			logger.EnableDebugMode()
		}

		overrideFilePath, _ := cmd.Flags().GetString("config")
		config_manager.Configure(overrideFilePath)
		initializeModuleWhichRequireConfig()
	},
}

func initializeModuleWhichRequireConfig() {
	notifictions.Initialize()
	services.InitializeRMQ()
	services.InitializeSQS()
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		//Remove the release-please comment from the version output
		version = strings.Split(version, "##")[0]
		fmt.Printf("Magecomm Version: %s\n", strings.TrimSpace(version))
		return
	}

	RootCmd.AddCommand(cmd.ListenCmd)
	RootCmd.AddCommand(cmd.MagerunCmd)
	RootCmd.AddCommand(cmd.CatCmd)
	RootCmd.AddCommand(cmd.OutputsCmd)

	RootCmd.PersistentFlags().String("config", "", "Path to config file")
	RootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")

	err := RootCmd.Execute()
	if err != nil {
		logger.Fatalf("Failed to execute command: %s", strings.ReplaceAll(err.Error(), "\n", " "))
	}
}
