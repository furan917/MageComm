package magerun

import (
	"fmt"
	"magecomm/config_manager"
	"magecomm/logger"
	"os"
	"os/exec"
	"strings"
)

const CommandMageRun = "magerun"

var defaultAllowedCommands = []string{
	"cache:clean",
	"admin:user:create",
	"admin:user:delete",
	"sys:cron:info",
	"setup:upgrade",
	"set:di:compile",
	"index:status",
	"index:reset",
	"index:reindex",
}

func HandleMagerunCommand(messageBody string) {
	command, args := parseMagerunCommand(messageBody)
	if !IsCommandAllowed(command) {
		return
	}
	//re-merge the command and args for execution
	args = append([]string{command}, args...)
	executeMagerunCommand(args)
}

func parseMagerunCommand(messageBody string) (string, []string) {
	args := strings.Fields(messageBody)
	return args[0], args[1:]
}

// IsCommandAllowed todo:: switch to a more secure method of checking allowed commands, envs can be overwritten by the user
func IsCommandAllowed(command string) bool {
	var allowedCommands []string

	allowedCommandsEnv := config_manager.GetValue(config_manager.CommandConfigAllowedMageRunCommands)
	if allowedCommandsEnv != "" {
		allowedCommands = strings.Split(allowedCommandsEnv, ",")
	} else {
		allowedCommands = defaultAllowedCommands
	}

	for _, allowedCommand := range allowedCommands {
		if allowedCommand == command {
			return true
		}
	}
	// print allowed commands
	logger.Warnf("Command not allowed, allowed commands are:\n%s \n", strings.Join(allowedCommands, ",\n"))
	fmt.Printf("%s Command not allowed, allowed commands are:\n%s \n", command, strings.Join(allowedCommands, ",\n"))
	return false
}

func executeMagerunCommand(args []string) {
	logger.Infof("Executing command %s with args: %v\n", CommandMageRun, args)
	cmd := exec.Command(CommandMageRun, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		logger.Warnf("Error executing magerun command:", err)
	}
}
