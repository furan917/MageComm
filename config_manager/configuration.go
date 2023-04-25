package config_manager

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"magecomm/logger"
	"os"
	"runtime"
	"strings"
)

const (
	CommandConfigLogPath                      = "magecomm_log_path"
	CommandConfigLogLevel                     = "magecomm_log_level"
	CommandConfigMaxOperationalCpuLimit       = "magecomm_max_operational_cpu_limit"
	CommandConfigMaxOperationalMemoryLimit    = "magecomm_max_operational_memory_limit"
	CommandConfigEnvironment                  = "magecomm_environment"
	CommandConfigListenerEngine               = "magecomm_listener_engine"
	CommandConfigListeners                    = "magecomm_listeners"
	CommandConfigPublisherOutputTimeout       = "magecomm_publisher_output_timeout"
	CommandConfigMageRunCommandPath           = "magecomm_magerun_command_path"
	CommandConfigAllowedMageRunCommands       = "magecomm_allowed_magerun_commands"
	CommandConfigRestrictedMagerunCommandArgs = "magecomm_restricted_magerun_command_args"
	CommandConfigRequiredMagerunCommandArgs   = "magecomm_required_magerun_command_args"
	CommandConfigDeployArchiveFolder          = "magecomm_deploy_archive_path"
	CommandConfigDeployArchiveLatestFile      = "magecomm_deploy_archive_latest_file"

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

func getDefault(key string) string {
	// we cant use viper.setDefault due to the order of operations we need: Config > Env > Default
	defaults := map[string]string{
		CommandConfigLogPath:                   "",
		CommandConfigLogLevel:                  "warn",
		CommandConfigMaxOperationalCpuLimit:    "80",
		CommandConfigMaxOperationalMemoryLimit: "80",
		CommandConfigEnvironment:               "default",
		CommandConfigListenerEngine:            "sqs",
		CommandConfigListeners:                 "",
		CommandConfigPublisherOutputTimeout:    "60",
		CommandConfigMageRunCommandPath:        "",
		CommandConfigAllowedMageRunCommands:    "",
		CommandConfigDeployArchiveFolder:       "/srv/magecomm/deploy/",
		CommandConfigDeployArchiveLatestFile:   "latest.tar.gz",
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

var defaultAllowedCommands = []string{
	"admin:token:create",
	"admin:user:unlock",
	"app:config:import",
	"braintree:migrate",
	"cache:clean",
	"cache:disable",
	"cache:enable",
	"cache:flush",
	"catalog:images:resize",
	"catalog:product:attributes:cleanup",
	"cms:block:toggle",
	"cms:wysiwyg:restrict",
	"cron:install",
	"cron:remove",
	"cron:run",
	"dev:query-log:disable",
	"dev:query-log:enable",
	"downloadable:domains:add",
	"downloadable:domains:remove",
	"inchoo:catalog:footwear-link-update",
	"index:trigger:recreate",
	"indexer:reindex",
	"indexer:reset",
	"indexer:set-mode",
	"klevu:images",
	"klevu:rating",
	"klevu:sync:category",
	"klevu:sync:cmspages",
	"klevu:syncdata",
	"klevu:syncstore:storecode",
	"maintenance:allow-ips",
	"maintenance:disable",
	"maintenance:enable",
	"media:dump",
	"msp:security:recaptcha:disable",
	"queue:consumers:start",
	"sys:cron:run",
	"sys:maintenance",
	"yotpo:order",
	"yotpo:sync",
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
			logger.Warnf("Failed to read the config file, reading from ENV vars, this is less secure: %v", err)
			return
		}
	}

	if logPath := GetValue(CommandConfigLogPath); logPath != "" {
		logger.ConfigureLogPath(logPath)
	}

	if logLevel := GetValue(CommandConfigLogLevel); logLevel != "" {
		logger.SetLogLevel(logLevel)
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

func IsMageRunCommandAllowed(command string) bool {
	var allowedCommands []string

	allowedCommandsEnv := GetValue(CommandConfigAllowedMageRunCommands)
	if allowedCommandsEnv != "" {
		allowedCommands = strings.Split(allowedCommandsEnv, ",")
	} else {
		allowedCommands = defaultAllowedCommands
	}

	for _, allowedCommand := range allowedCommands {
		if allowedCommand == command {
			return true
		}
	}
	// print allowed commands
	logger.Warnf("Command not allowed, allowed commands are:\n%s \n", strings.Join(allowedCommands, ",\n"))
	fmt.Printf("%s Command not allowed, allowed commands are:\n%s \n", command, strings.Join(allowedCommands, ",\n"))
	return false
}

func IsRestrictedCommandArgsIncluded(command string, args []string) bool {
	restrictedCommandArgMap := ParseCommandArgsMap(GetValue(CommandConfigRestrictedMagerunCommandArgs))

	restrictedArgsList, commandExists := restrictedCommandArgMap[command]
	if !commandExists {
		return false
	}
	// in go it is more performant to use maps with null/"" values to reduce search complexity from a linear (O(n)) to a constant (O(1))
	// but for the sake of configuration simplicity, we use a mapped list
	for _, arg := range args {
		for _, restrictedArg := range restrictedArgsList {
			if arg == restrictedArg {
				return true
			}
		}
	}

	return false
}

func IsRequiredCommandArgsIncluded(command string, args []string) (bool, []string) {
	requiredCommandArgMap := ParseCommandArgsMap(GetValue(CommandConfigRequiredMagerunCommandArgs))
	requiredArgsList, commandExists := requiredCommandArgMap[command]
	if !commandExists {
		return true, []string{}
	}

	for i := 0; i < len(requiredArgsList); i++ {
		requiredArg := requiredArgsList[i]
		for _, arg := range args {
			if arg == requiredArg {
				requiredArgsList = append(requiredArgsList[:i], requiredArgsList[i+1:]...)
				i-- // adjust the index since we removed an element
				break
			}
		}
	}

	if len(requiredArgsList) == 0 {
		return true, []string{}
	}

	return false, requiredArgsList
}
