package magerun

import (
	"bytes"
	"fmt"
	"magecomm/config_manager"
	"magecomm/logger"
	"os/exec"
	"strings"
)

const DefaultCommandMageRun = "magerun"

func HandleMagerunCommand(messageBody string) (string, error) {
	command, args := parseMagerunCommand(messageBody)
	if !config_manager.IsMageRunCommandAllowed(command) {
		return "", fmt.Errorf("command %s is not allowed", command)
	}

	if config_manager.IsRestrictedCommandArgsIncluded(command, args) {
		return "", fmt.Errorf("the command '%s' is not allowed with the following arguments: %s", command, strings.Join(args, " "))
	}

	if isAllRequiredArgsIncluded, missingRequiredArgs := config_manager.IsRequiredCommandArgsIncluded(command, args); !isAllRequiredArgsIncluded {
		return "", fmt.Errorf("the command '%s' is missing some required arguments: %s, unable to run command", command, strings.Join(missingRequiredArgs, " "))
	}

	args = append([]string{command}, args...)
	return executeMagerunCommand(args)
}

func executeMagerunCommand(args []string) (string, error) {
	mageRunCmdPath := getMageRunCommand()
	logger.Infof("Executing command %s with args: %v\n", mageRunCmdPath, args)
	cmd := exec.Command(mageRunCmdPath, args...)

	var stdoutBuffer, stderrBuffer bytes.Buffer
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stderrBuffer

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error executing magerun command: %s", err)
	}

	stdoutStr := stdoutBuffer.String()
	stderrStr := stderrBuffer.String()

	output := stdoutStr + "\n" + stderrStr
	return output, nil
}

func getMageRunCommand() string {
	configuredMageRunCmd := config_manager.GetValue(config_manager.CommandConfigMageRunCommandPath)
	if configuredMageRunCmd == "" {
		return DefaultCommandMageRun
	}

	return configuredMageRunCmd
}

func parseMagerunCommand(messageBody string) (string, []string) {
	args := strings.Fields(messageBody)
	return args[0], args[1:]
}
