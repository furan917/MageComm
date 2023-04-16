package magerun

import (
	"bytes"
	"fmt"
	"magecomm/config_manager"
	"magecomm/logger"
	"magecomm/messages/publisher"
	"magecomm/messages/queues"
	"os/exec"
	"strings"
)

const CommandMageRun = "magerun"

func parseMagerunCommand(messageBody string) (string, []string) {
	args := strings.Fields(messageBody)
	return args[0], args[1:]
}

func HandleMagerunCommand(messageBody string, correlationID string) {
	command, args := parseMagerunCommand(messageBody)
	if !config_manager.IsMageRunCommandAllowed(command) {
		return
	}
	args = append([]string{command}, args...)
	output, err := executeMagerunCommand(args)

	// Publish the output to the RMQ/SQS queue
	publisher, err := publisher.MapPublisherToEngine()
	if err != nil {
		logger.Warnf("Error publishing message to RMQ/SQS queue:", err)
	}

	_, err = publisher.PublishMessage(output, queues.MapQueueToOutputQueue(CommandMageRun), correlationID)
	if err != nil {
		logger.Errorf("failed to publish message: %v", err)
	}
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
