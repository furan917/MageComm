package config_manager

import (
	"github.com/spf13/viper"
	"magecomm/logger"
	"os"
	"runtime"
	"strings"
)

const (
	CommandConfigEnvironment            = "environment"
	CommandConfigListenerEngine         = "listener_engine"
	CommandConfigListeners              = "listeners"
	CommandConfigAllowedMageRunCommands = "allowed_magerun_commands"
)

func getDefault(key string) string {
	// we cant use viper.setDefault due to the order of operations we need: Config > Env > Default
	defaults := map[string]string{
		CommandConfigEnvironment:            "default",
		CommandConfigListenerEngine:         "sqs",
		CommandConfigListeners:              "",
		CommandConfigAllowedMageRunCommands: "",
		"sqs_aws_region":                    "eu-west-1",
		"rmq_host":                          "localhost",
		"rmq_port":                          "5672",
		"rmq_user":                          "guest",
		"rmq_pass":                          "guest",
		"rmq_vhost":                         "/",
		"rmq_tls":                           "false",
	}

	if defaultValue, ok := defaults[key]; ok {
		return defaultValue
	}

	return ""
}

func Configure() {
	viper.SetConfigName("config")
	if runtime.GOOS == "windows" {
		viper.AddConfigPath(os.Getenv("APPDATA") + "\\magecomm\\")
	} else {
		viper.AddConfigPath("/etc/magecomm/")
	}
	err := viper.ReadInConfig()
	if err != nil {
		// If the configuration file does not exist, warn user that env vars will be used
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Infof("No config file found, reading fully from env vars, this is less secure")
		} else {
			logger.Warnf("Failed to read the config file, reading from ENV vars, this is less secure:", err)
			return
		}
	}
}

func GetValue(key string) string {
	if value, ok := getConfigValue(strings.ToLower(key)); ok {
		return value
	}
	if value, ok := getEnvFallback(strings.ToUpper(key)); ok {
		return value
	}
	value := getDefault(strings.ToLower(key))
	if value == "" {
		logger.Debugf("No config, env, or default value set for  %s", key)
	}

	return value
}

func getConfigValue(key string) (string, bool) {
	value := viper.GetString(key)
	if value != "" {
		return value, true
	}
	return "", false
}

func getEnvFallback(key string) (string, bool) {
	logger.Infof("No config value set for key %s, using fallback ENV, this is less secure as users can manipulate these values", key)
	value, ok := os.LookupEnv(key)
	if ok && value != "" {
		return value, true
	}
	return "", false
}
