package config_manager

import (
    "encoding/json"
    "github.com/spf13/viper"
    "magecomm/logger"
    "os"
    "runtime"
    "strings"
)

const (
    ConfigLogPath                             = "magecomm_log_path"
    ConfigLogLevel                            = "magecomm_log_level"
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

    // Slack
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

var trueValues = []string{"true", "1", "yes", "y"}

func getDefault(key string) string {
    // we cant use viper.setDefault due to the order of operations we need: Config > Env > Default
    defaults := map[string]string{
        ConfigLogPath:                          "",
        ConfigLogLevel:                         "warn",
        CommandConfigMaxOperationalCpuLimit:    "80",
        CommandConfigMaxOperationalMemoryLimit: "80",
        CommandConfigEnvironment:               "default",
        CommandConfigListenerEngine:            "sqs",
        CommandConfigListeners:                 "",
        CommandConfigPublisherOutputTimeout:    "600",
        CommandConfigMageRunCommandPath:        "",
        CommandConfigAllowedMageRunCommands:    "",
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
        viper.SetConfigName("config")
        if runtime.GOOS == "windows" {
            viper.AddConfigPath(os.Getenv("APPDATA") + "\\magecomm\\")
        } else {
            viper.AddConfigPath("/etc/magecomm/")
        }
    }

    configName := viper.ConfigFileUsed()
    logger.Infof("Using config file: %s", configName)
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
    }

    if logLevel := GetValue(ConfigLogLevel); logLevel != "" {
        logger.SetLogLevel(logLevel)
    }
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
