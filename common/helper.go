package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Contains Helper function to check if a slice Contains a string
func Contains(slice []string, item string) bool {
	for _, substring := range slice {
		if strings.Contains(item, substring) {
			return true
		}
	}
	return false
}

func PromptUserForConfirmation(prompt string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s [y/N]: ", prompt)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	response = strings.ToLower(strings.TrimSpace(response))
	if response == "y" || response == "yes" {
		return true, nil
	}
	return false, nil
}
