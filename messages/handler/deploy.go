package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"magecomm/config_manager"
	"magecomm/logger"
	"magecomm/magerun"
	"magecomm/messages/publisher"
	"magecomm/messages/queues"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type DeployHandler struct {
	Publisher publisher.MessagePublisher
}

//todo:: convert GetValues to target CONSTs

func (handler *DeployHandler) ProcessMessage(messageBody string, correlationID string) error {

	// Execute pre-hooks
	preHooksValue := config_manager.GetValue("pre_hooks")
	preHooks := strings.Split(preHooksValue, ",")
	for _, hook := range preHooks {
		fmt.Printf("Running pre-hook: %s\n", hook)
		err := handler.executeScript(hook)
		if err != nil {
			return fmt.Errorf("failed to execute pre-hook %s: %v", hook, err)
		}
	}

	// Fetch and unpack code
	err := handler.fetchAndUnpack(messageBody)
	if err != nil {
		return fmt.Errorf("failed to fetch and unpack code: %v", err)
	}

	// Update symlinks
	err = handler.updateSymlinks()
	if err != nil {
		return fmt.Errorf("failed to update symlinks: %v", err)
	}

	var steps [][]string
	var isBlueGreenEnabled = handler.isBlueGreenEnabled()
	if isBlueGreenEnabled {
		steps = append(steps,
			[]string{"maintenance:enable"},
		)
	}

	if handler.hasDBChanges() {
		steps = append(steps, []string{"setup:upgrade", "--keep-generated"})
	}
	steps = append(steps,
		[]string{"cache:clean"},
	)

	if isBlueGreenEnabled {
		steps = append(steps,
			[]string{"maintenance:enable"},
		)
	}

	for _, step := range steps {
		if err := handler.executeMagentoStep(step, correlationID); err != nil {
			return err
		}
	}

	postHooksValue := config_manager.GetValue("post_hooks")
	postHooks := strings.Split(postHooksValue, ",")
	for _, hook := range postHooks {
		fmt.Printf("Running post-hook: %s\n", hook)
		err := handler.executeScript(hook)
		if err != nil {
			return fmt.Errorf("failed to execute post-hook %s: %v", hook, err)
		}
	}

	fmt.Println("Deployment completed successfully.")
	return nil
}

func (handler *DeployHandler) executeMagentoStep(args []string, correlationID string) error {
	output, err := handler.runMagentoCommand(args)
	if err != nil {
		output += err.Error()
	}
	if err := handler.publishOutput(output, correlationID); err != nil {
		logger.Warnf("Error publishing message to RMQ/SQS queue:", err)
	}
	return err
}

func (handler *DeployHandler) isBlueGreenEnabled() bool {
	value := config_manager.GetValue("blue_green_enabled")
	if value != "1" && strings.ToLower(value) != "true" && strings.ToLower(value) != "yes" {
		return true
	}
	return false
}

func (handler *DeployHandler) hasDBChanges() bool {
	output, err := handler.runMagentoCommand([]string{"setup:db:status"})
	if err != nil {
		return true
	}
	if output == "All modules are up to date." {
		return false
	}

	return true
}

func (handler *DeployHandler) runMagentoCommand(args []string) (string, error) {
	output, err := magerun.ExecuteMagerunCommand(args)
	if err != nil {
		return output, fmt.Errorf("failed to run Magento command %v: %v\noutput:\n%s", args, err, output)
	}

	return output, nil
}

func (handler *DeployHandler) publishOutput(output string, correlationID string) error {
	publisherClass, err := publisher.MapPublisherToEngine()
	if err != nil {
		logger.Warnf("Error publishing message to RMQ/SQS queue:", err)
	}
	logger.Debugf("Publishing output to queue:", queues.MapQueueToOutputQueue("deploy"), "with correlation ID:", correlationID)
	_, err = publisherClass.Publish(output, queues.MapQueueToOutputQueue("deploy"), correlationID)
	if err != nil {
		return err
	}

	return nil
}

func (handler *DeployHandler) fetchAndUnpack(fileName string) error {
	logger.Debugf("Fetching and unpacking file:", fileName)
	//todo:: update fetch and Uppack to use S3/Github/Bitbucket/Archive on filesystem
	source := config_manager.GetValue("deploy_source")
	resp, err := http.Get(source)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Warnf("Error closing http response body: %v", err)
		}
	}(resp.Body)

	targetDir := config_manager.GetValue("target_dir")
	out, err := os.Create(targetDir)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			logger.Warnf("Error closing file: %v", err)
		}
	}(out)

	_, err = io.Copy(out, resp.Body)
	return err
}

func (handler *DeployHandler) updateSymlinks() error {
	symlinks := config_manager.GetValue("symlinks")
	var symlinksJson map[string]string
	err := json.Unmarshal([]byte(symlinks), &symlinksJson)
	if err != nil {
		return fmt.Errorf("failed to convert symlinks to json: %v", err)
	}
	for src, dest := range symlinksJson {
		err := os.Symlink(src, dest)
		if err != nil {
			return fmt.Errorf("failed to create symlink %s -> %s: %v", src, dest, err)
		}
	}

	targetDir := config_manager.GetValue("target_dir")
	documentRoot := config_manager.GetValue("document_root")
	err = os.Symlink(targetDir, documentRoot)
	if err != nil {
		return fmt.Errorf("failed to update document root symlink: %v", err)
	}

	return nil
}

func (handler *DeployHandler) executeScript(scriptPath string) error {
	scriptAbsPath, err := filepath.Abs(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of script %s: %v", scriptPath, err)
	}

	cmd := exec.Command(scriptAbsPath)
	targetDir := config_manager.GetValue("target_dir")
	cmd.Dir = targetDir
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to execute script %s: %v\noutput:\n%s", scriptPath, err, string(output))
	}

	fmt.Printf("Script executed successfully: %s\noutput:\n%s", scriptPath, string(output))
	return nil
}
