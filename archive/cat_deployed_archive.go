package archive

import (
	"errors"
	"fmt"
	"magecomm/config_manager"
	"magecomm/logger"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func isSupportedArchive(filename string) bool {
	for _, ext := range SupportedCatArchives {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return true
		}
	}
	return false
}

func GetLatestDeploy() (string, error) {
	deployPath := config_manager.GetValue(config_manager.CommandConfigDeployArchiveFolder)
	if deployPath == "" {
		return "", errors.New("deploy path not set, please contact your system administrator")
	}

	if _, err := os.Stat(deployPath); os.IsNotExist(err) {
		return "", errors.New("deploy path set but does not exist, please contact your system administrator")
	}

	fileName, err := getLatestArchiveFileName(deployPath)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func getLatestArchiveFileName(deployPath string) (string, error) {
	_, err := CheckFolder(deployPath)
	if err != nil {
		return "", fmt.Errorf("failed to check deploy folder: %w", err)
	} else {
		logger.Debugf("Deploy path exists: %s", deployPath)
	}
	fileName := config_manager.GetValue(config_manager.CommandConfigDeployArchiveLatestFile)
	if fileName != "" {
		if _, err := os.Stat(filepath.Join(deployPath, fileName)); err == nil {
			return filepath.Join(deployPath, fileName), nil
		}
	}

	//configured deploy file not found, loop through all files that are archives in the deployment folder and return the latest one
	logger.Debugf("No configured deploy file found, looping through all files in deploy folder")
	files, err := os.ReadDir(deployPath)
	if err != nil {
		return "", err
	}

	sort.Sort(sort.Reverse(ByModTime(files)))
	for _, file := range files {
		if !file.IsDir() && isSupportedArchive(file.Name()) {
			return filepath.Join(deployPath, file.Name()), nil
		}
	}

	return "", errors.New("no supported archive file found")
}

func CheckFolder(deployPath string) (string, error) {
	_, err := os.Stat(deployPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New("deploy path set but does not exist, please contact your system administrator")
		} else if os.IsPermission(err) {
			return "", errors.New("deploy path exists, but no permission to read, please contact your system administrator")
		} else {
			return "", fmt.Errorf("unknown error while checking deploy path: %v", err)
		}
	}
	return deployPath, nil
}
