package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"magecomm/common"
	"magecomm/config_manager"
	"magecomm/logger"
	"magecomm/messages/listener"
	"magecomm/messages/publisher"
	"magecomm/notifictions"
	"magecomm/services"
	"strings"
)

const MageRunQueue = "magerun"

var MagerunCmd = &cobra.Command{
	Use:                "magerun",
	Short:              "A wrapper for the magerun command with restricted command usage",
	DisableFlagParsing: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		//empty pre run to stop execution of parent RootCmd pre run
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// handle global arguments (e.g. --config, --debug) as root command cannot due to DisableFlagParsing
		var globalArguments = handleGlobalArguments(args)
		magerunArgs := args[len(globalArguments):]
		if len(magerunArgs) < 1 {
			return fmt.Errorf("no command provided")
		}

		command := magerunArgs[0]
		if isCmdAllowed, err := config_manager.IsMageRunCommandAllowed(command); !isCmdAllowed {
			return err
		}

		if isRestrictedArgsIncluded, err := config_manager.IsRestrictedCommandArgsIncluded(command, magerunArgs[1:]); isRestrictedArgsIncluded {
			return err
		}

		if isAllRequiredArgsIncluded, missingRequiredArgs := config_manager.IsRequiredCommandArgsIncluded(command, magerunArgs[1:]); !isAllRequiredArgsIncluded {
			prompt := fmt.Sprintf("The command '%s' is missing required arguments: %s. Do you want to run this command and include them?", command, strings.Join(missingRequiredArgs, " "))
			confirmed, err := common.PromptUserForConfirmation(prompt)
			if err != nil {
				return fmt.Errorf("error while reading user input: %v", err)
			}

			if confirmed {
				magerunArgs = append(magerunArgs, missingRequiredArgs...)
			} else {
				return fmt.Errorf("exiting: the command '%s' cannot be executed without the required arguments: %s", command, strings.Join(missingRequiredArgs, " "))
			}
		}

		err := handleMageRunCmdMessage(magerunArgs)
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
		logger.Warnf("Command executed, but no output could be returned")
		fmt.Println("Command executed, but no output could be returned")
		return nil
	}

	output, err := listener.HandleOutputByCorrelationID(correlationID, MageRunQueue)
	if err != nil {
		return fmt.Errorf("failed to get output: %s", err)
	}

	if output != "" {
		logger.Infof("Output printed to terminal")
		fmt.Println(output)
	}

	return nil
}

func initializeModuleWhichRequireConfig() {
	notifictions.Initialize()
	services.InitializeRMQ()
	services.InitializeSQS()
}

func handleGlobalArguments(args []string) []string {
	// Replicates RootCmd.PersistentPreRunE as it is not usable when DisableFlagParsing is set to true
	// global Arg handling must be done manually
	var globalArguments []string
	var overrideFilePath string

	//Note arguments between start and magerun command.
	for i, arg := range args {
		if strings.HasPrefix(arg, "--") {
			globalArguments = append(globalArguments, arg)

			//if arg is --debug then enable debug mode
			if arg == "--debug" {
				logger.Warnf("Magerun: Debug mode enabled")
				logger.EnableDebugMode()
			}

			//if arg is --config then configure
			if strings.HasPrefix(arg, "--config") {
				var configPath string
				if strings.Contains(arg, "=") {
					configPath = strings.Split(arg, "=")[1]
				} else {
					//ensure next argument is config path
					if len(args) > i+1 && !strings.HasPrefix(args[i+1], "--") {
						configPath = args[i+1]
						//Ensure we remove the config path from args to avoid breaking early
						args = append(args[:i+1], args[i+2:]...)
					}
				}
				if configPath == "" {
					logger.Warnf("No config file path provided with argument, using default config file path")
				}

				overrideFilePath = strings.Trim(configPath, `"'`)
			}
		} else {
			//We have reached the magerun command, so exit loop
			break
		}
	}

	config_manager.Configure(overrideFilePath)
	initializeModuleWhichRequireConfig()

	return globalArguments
}
