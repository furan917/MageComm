package common

import "strings"

// Contains Helper function to check if a slice Contains a string
func Contains(slice []string, item string) bool {
	for _, substring := range slice {
		if strings.Contains(item, substring) {
			return true
		}
	}
	return false
}
