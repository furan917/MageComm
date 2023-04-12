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
	"admin:token:create",
	"admin:user:unlock",
	"app:config:import",
	"braintree:migrate",
	"cache:clean",
	"cache:disable",
	"cache:enable",
	"cache:flush",
	"catalog:images:resize",
	"catalog:product:attributes:cleanup",
	"cms:block:toggle",
	"cms:wysiwyg:restrict",
	"cron:install",
	"cron:remove",
	"cron:run",
	"dev:query-log:disable",
	"dev:query-log:enable",
	"downloadable:domains:add",
	"downloadable:domains:remove",
	"inchoo:catalog:footwear-link-update",
	"index:trigger:recreate",
	"indexer:reindex",
	"indexer:reset",
	"indexer:set-mode",
	"klevu:images",
	"klevu:rating",
	"klevu:sync:category",
	"klevu:sync:cmspages",
	"klevu:syncdata",
	"klevu:syncstore:storecode",
	"maintenance:allow-ips",
	"maintenance:disable",
	"maintenance:enable",
	"media:dump",
	"msp:security:recaptcha:disable",
	"queue:consumers:start",
	"sys:cron:run",
	"sys:maintenance",
	"yotpo:order",
	"yotpo:sync",
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
