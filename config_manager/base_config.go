package config_manager

import (
	"encoding/json"
	"github.com/spf13/viper"
	"magecomm/common"
	"magecomm/logger"
	"os"
	"runtime"
	"strings"
)

const (
	ConfigLogPath                             = "magecomm_log_path"
	ConfigLogLevel                            = "magecomm_log_level"
	ConfigPrintoutLogLevel                    = "magecomm_printout_log_level"
	CommandConfigMaxOperationalCpuLimit       = "magecomm_max_operational_cpu_limit"
	CommandConfigMaxOperationalMemoryLimit    = "magecomm_max_operational_memory_limit"
	CommandConfigEnvironment                  = "magecomm_environment"
	CommandConfigListenerEngine               = "magecomm_listener_engine"
	CommandConfigListeners                    = "magecomm_listeners"
	CommandConfigAllowedQueues                = "magecomm_listener_allowed_queues"
	CommandConfigPublisherOutputTimeout       = "magecomm_publisher_output_timeout"
	CommandConfigMageRunCommandPath           = "magecomm_magerun_command_path"
	CommandConfigAllowedMageRunCommands       = "magecomm_allowed_magerun_commands"
	CommandConfigRestrictedMagerunCommandArgs = "magecomm_restricted_magerun_command_args"
	CommandConfigRequiredMagerunCommandArgs   = "magecomm_required_magerun_command_args"
	CommandConfigForceMagerunNoInteraction    = "magecomm_force_magerun_no_interaction"
	CommandConfigDeployArchiveFolder          = "magecomm_deploy_archive_path"
	CommandConfigDeployArchiveLatestFile      = "magecomm_deploy_archive_latest_file"

	//Slack
	ConfigSlackEnabled                    = "magecomm_slack_enabled"
	ConfigSlackDisableOutputNotifications = "magecomm_slack_disable_output_notifications"

	ConfigSlackWebhookUrl      = "magecomm_slack_webhook_url"
	ConfigSlackWebhookChannel  = "magecomm_slack_webhook_channel"
	ConfigSlackWebhookUserName = "magecomm_slack_webhook_username"
	ConfigSlackAppToken        = "magecomm_slack_app_token"
	ConfigSlackAppChannel      = "magecomm_slack_app_channel"
	ConfigSlackAppUserName     = "magecomm_slack_app_username"

	//SQS
	ConfigSQSRegion = "magecomm_sqs_aws_region"

	//RMQ
	ConfigRabbitMQTLS   = "magecomm_rmq_tls"
	ConfigRabbitMQUser  = "magecomm_rmq_user"
	ConfigRabbitMQPass  = "magecomm_rmq_pass"
	ConfigRabbitMQHost  = "magecomm_rmq_host"
	ConfigRabbitMQPort  = "magecomm_rmq_port"
	ConfigRabbitMQVhost = "magecomm_rmq_vhost"
)

// Slice values e.g Key => [1,2,3]
var sliceValues = []string{
	CommandConfigListeners,
	CommandConfigAllowedQueues,
	CommandConfigAllowedMageRunCommands,
}

// Map values e.g MappedKey => { Key1 => [1,2,3], Key2 => [4,5,6] }
var mappedValues = []string{
	CommandConfigRestrictedMagerunCommandArgs,
	CommandConfigRequiredMagerunCommandArgs,
}

var trueValues = []string{"true", "1", "yes", "y"}

func getDefault(key string) string {
	// we cant use viper.setDefault due to the order of operations we need: Config > Env > Default
	defaults := map[string]string{
		ConfigLogPath:                          "",
		ConfigLogLevel:                         "warn",
		ConfigPrintoutLogLevel:                 "error",
		CommandConfigMaxOperationalCpuLimit:    "80",
		CommandConfigMaxOperationalMemoryLimit: "80",
		CommandConfigEnvironment:               "default",
		CommandConfigListenerEngine:            "sqs",
		CommandConfigAllowedQueues:             "cat,magerun",
		CommandConfigListeners:                 "",
		CommandConfigPublisherOutputTimeout:    "600",
		CommandConfigMageRunCommandPath:        "",
		CommandConfigAllowedMageRunCommands:    "",
		CommandConfigForceMagerunNoInteraction: "true",
		CommandConfigDeployArchiveFolder:       "/srv/magecomm/deploy/",
		CommandConfigDeployArchiveLatestFile:   "latest.tar.gz",
		ConfigSlackEnabled:                     "false",
		ConfigSlackDisableOutputNotifications:  "false",
		ConfigSQSRegion:                        "eu-west-1",
		ConfigRabbitMQTLS:                      "false",
		ConfigRabbitMQUser:                     "guest",
		ConfigRabbitMQPass:                     "guest",
		ConfigRabbitMQHost:                     "localhost",
		ConfigRabbitMQPort:                     "5672",
		ConfigRabbitMQVhost:                    "/",
	}

	if defaultValue, ok := defaults[key]; ok {
		return defaultValue
	}

	return ""
}

func Configure(overrideFile string) {
	if overrideFile != "" {
		viper.SetConfigFile(overrideFile)
	} else {
		// Set base folder and file name of config file
		viper.SetConfigName("config")
		if runtime.GOOS == "windows" {
			viper.AddConfigPath(os.Getenv("APPDATA") + "\\magecomm\\")
		} else {
			viper.AddConfigPath("/etc/magecomm/")
		}
		// Search for json config file first, then fallback to yaml
		viper.SetConfigType("json")
		if err := viper.ReadInConfig(); err != nil {
			viper.SetConfigType("yaml")
		}
	}

	err := viper.ReadInConfig()
	if err != nil {
		// If the configuration file does not exist, warn user that env vars will be used
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Infof("No config file found, reading fully from env vars, this is less secure")
		} else {
			logger.Warnf("Failed to read the config file, reading from ENV vars, this is less secure: %v", err)
			return
		}
	}

	if logPath := GetValue(ConfigLogPath); logPath != "" {
		logger.ConfigureLogPath(logPath)
		logger.Infof("Logging to file: %s", logPath)
	}

	if logLevel := GetValue(ConfigLogLevel); logLevel != "" {
		logger.SetLogLevel(logLevel)
		logger.Infof("Logging level set to: %s", logLevel)
	}

	configName := viper.ConfigFileUsed()
	logger.Infof("Using config file: %s", configName)
}

func GetBoolValue(key string) bool {
	value := GetValue(key)
	for _, v := range trueValues {
		if strings.ToLower(value) == v {
			return true
		}
	}
	return false
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

	// if map, slice or string use different methods
	if common.Contains(sliceValues, key) {
		value := viper.GetStringSlice(key)
		if len(value) > 0 {
			return strings.Join(value, ","), true
		}
	} else if common.Contains(mappedValues, key) {
		value := viper.GetStringMapStringSlice(key)
		if len(value) > 0 {
			jsonString, _ := json.Marshal(value)
			return string(jsonString), true
		}
	} else {
		value := viper.GetString(key)
		if value != "" {
			return value, true
		}
	}

	return "", false
}

func getEnvFallback(key string) (string, bool) {
	value, ok := os.LookupEnv(key)
	if ok && value != "" {
		logger.Infof("No config value set for key %s, using fallback ENV, this is less secure as users can manipulate these values", key)

		return value, true
	}
	return "", false
}

func ParseCommandArgsMap(jsonString string) map[string][]string {
	var commandArgsMap map[string][]string
	if err := json.Unmarshal([]byte(jsonString), &commandArgsMap); err != nil {
		logger.Warnf("Failed to parse passed in command args JSON: %s", err)
		return map[string][]string{}
	}
	return commandArgsMap
}

func GetEngine() string {
	return GetValue(CommandConfigListenerEngine)
}
