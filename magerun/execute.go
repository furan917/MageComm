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

	//if --no-interaction/-n is not set, set it
	forceNoInteractionFlag := config_manager.GetBoolValue(config_manager.CommandConfigForceMagerunNoInteraction)
	noInteractionFlagPresent := false
	for _, arg := range args {
		if arg == "--no-interaction" || arg == "-n" {
			noInteractionFlagPresent = true
			break
		}
	}
	if !noInteractionFlagPresent && forceNoInteractionFlag {
		logger.Infof("The command '%s' does not contain the '--no-interaction' flag, adding it to the command", command)
		args = append(args, "--no-interaction")
	}

	args = append([]string{command}, args...)
	return executeMagerunCommand(args)
}

func executeMagerunCommand(args []string) (string, error) {
	mageRunCmdPath := getMageRunCommand()
	logger.Infof("Executing command %s with args: %v", mageRunCmdPath, args)

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
	// Grab any output before returning with command error
	stdoutStr := stdoutBuffer.String()
	stderrStr := stderrBuffer.String()
	output := stripMagerunOutput(stdoutStr + "\n" + stderrStr)

	// Now check command for error and return either success or failure
	if err != nil {
		logger.Warnf("Error executing magerun command: %s, with the following output: %s", err, strings.ReplaceAll(output, "\n", " "))
		return output, fmt.Errorf("error executing magerun command: %s", err)
	}
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
	//trim any leading or trailing whitespace
	strippedOutput = strings.TrimSpace(strippedOutput)

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
	//if messageBody contains \" or \' then replace with " or '
	escapedQuotePattern := `\\(["'])`
	re := regexp.MustCompile(escapedQuotePattern)
	messageBody = re.ReplaceAllString(messageBody, `$1`)

	args := strings.Fields(messageBody)
	return args[0], args[1:]
}

// We absolutely must not allow command escaping. e.g. magerun cache:clean; rm -rf /
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
