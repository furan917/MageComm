package magerun

import (
	"bytes"
	"fmt"
	"magecomm/common"
	"magecomm/config_manager"
	"magecomm/logger"
	"magecomm/notifictions"
	"os/exec"
	"regexp"
	"strings"
)

const DefaultCommandMageRun = "magerun"

func HandleMagerunCommand(messageBody string) (string, error) {
	command, args := parseMagerunCommand(messageBody)
	args = sanitizeCommandArgs(args)

	if isCmdAllowed, err := config_manager.IsMageRunCommandAllowed(command); !isCmdAllowed {
		return "", err
	}

	if isRestrictedArgsIncluded, err := config_manager.IsRestrictedCommandArgsIncluded(command, args); isRestrictedArgsIncluded {
		return "", err
	}

	if isAllRequiredArgsIncluded, missingRequiredArgs := config_manager.IsRequiredCommandArgsIncluded(command, args); !isAllRequiredArgsIncluded {
		return "", fmt.Errorf("the command '%s' is missing some required arguments: %s, unable to run command", command, strings.Join(missingRequiredArgs, " "))
	}

	//if --no-interaction is not set, set it
	if !common.Contains(args, "--no-interaction") {
		args = append(args, "--no-interaction")
		logger.Infof("The command '%s' does not contain the '--no-interaction' flag, adding it to the command", command)
	}

	args = append([]string{command}, args...)
	return executeMagerunCommand(args)
}

func executeMagerunCommand(args []string) (string, error) {
	mageRunCmdPath := getMageRunCommand()
	logger.Infof("Executing command %s with args: %v\n", mageRunCmdPath, args)

	if config_manager.GetBoolValue(config_manager.ConfigSlackEnabled) {
		logger.Infof("Slack notification is enabled, sending notification")
		notifier := notifictions.DefaultSlackNotifier
		err := notifier.Notify(
			fmt.Sprintf("Executing command: '%v' on environment: '%s'", strings.Join(args, " "), config_manager.GetValue(config_manager.CommandConfigEnvironment)))
		if err != nil {
			logger.Warnf("Failed to send slack notification: %v", err)
		}
	}

	splitCmd := strings.Fields(mageRunCmdPath)
	cmd := exec.Command(splitCmd[0], splitCmd[1:]...)
	cmd.Args = append(cmd.Args, args...)

	var stdoutBuffer, stderrBuffer bytes.Buffer
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stderrBuffer

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error executing magerun command: %s", err)
	}

	stdoutStr := stdoutBuffer.String()
	stderrStr := stderrBuffer.String()

	output := stripMagerunOutput(stdoutStr + "\n" + stderrStr)

	logger.Infof("Executed command %s with args: %v and handling output", mageRunCmdPath, args)

	return output, nil
}

func stripMagerunOutput(output string) string {
	patterns := map[string]string{
		`(?i)(?:it's|it is) not recommended to run .*? as root user`: "",
		//Add more regex patterns here with their corresponding replacement
	}

	strippedOutput := output
	for pattern, replacement := range patterns {
		re := regexp.MustCompile(pattern)
		strippedOutput = re.ReplaceAllString(strippedOutput, replacement)
	}

	return strippedOutput
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

// We absolutely must not allow command escaping. e.g magerun cache:clean; rm -rf /
func sanitizeCommandArgs(args []string) []string {
	var sanitizedArgs []string
	disallowed := []string{";", "&&", "||", "|", "`", "$", "(", ")", "<", ">", "!"}
	for _, arg := range args {
		if common.Contains(disallowed, arg) {
			logger.Warnf("Command args contain potentially unsafe characters, removing arg: %s", arg)
			continue
		}
		sanitizedArgs = append(sanitizedArgs, arg)
	}
	return sanitizedArgs
}
