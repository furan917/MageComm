package magerun

import (
	"bytes"
	"fmt"
	"magecomm/config_manager"
	"magecomm/logger"
	"os/exec"
	"strings"
)

const CommandMageRun = "magerun"

func parseMagerunCommand(messageBody string) (string, []string) {
	args := strings.Fields(messageBody)
	return args[0], args[1:]
}

func HandleMagerunCommand(messageBody string) (string, error) {
	command, args := parseMagerunCommand(messageBody)
	if !config_manager.IsMageRunCommandAllowed(command) {
		return "", fmt.Errorf("command %s is not allowed", command)
	}
	args = append([]string{command}, args...)
	return executeMagerunCommand(args)
}

func executeMagerunCommand(args []string) (string, error) {
	logger.Infof("Executing command %s with args: %v\n", CommandMageRun, args)
	cmd := exec.Command(CommandMageRun, args...)

	var stdoutBuffer, stderrBuffer bytes.Buffer
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stderrBuffer

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error executing magerun command: %s", err)
	}

	stdoutStr := stdoutBuffer.String()
	stderrStr := stderrBuffer.String()

	// Combine stdout and stderr strings
	output := stdoutStr + "\n" + stderrStr
	return output, nil
}
