package config_manager

import "strings"

func IsAllowedQueue(queueName string) bool {
	if queueName == "" {
		return false
	}
	for _, allowedQueue := range GetAllowedQueues() {
		if queueName == allowedQueue {
			return true
		}
	}
	return false
}

func GetAllowedQueues() []string {
	return strings.Split(strings.ReplaceAll(GetValue(CommandConfigAllowedQueues), " ", ""), ",")
}
