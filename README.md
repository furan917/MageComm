# MageComm CLI Tool

MageComm CLI is a command line tool for managing Magento applications. It provides a convenient way to execute restricted magerun commands, manage deployments, and cat specific files within an archive.
The main use case for this tool is to provide a way to execute magerun commands in a controlled manner via a messaging service. This allows us to execute commands on a remote server without exposing the application server itself

*It is important to note that the environment configuration/env is important to set if you plan to use this tool in a shared rmq/sqs instance as this will prefix your queues to avoid cross communication*

The tool looks for a configuration file in `/etc/magecomm/`(unix) / `%APPDATA%\magecomm\`(windows) `config.yml|config.json` and if none found defaults to environment variables, then fallback to default values.


## Akoova Specific fork details
**Currently there is no custom code changes for akoova, this is a fork of the original repo to allow us to make changes if needed in the future**

- Release-Please is disabled for this repository.

Quick install for M1/M2 mac:
`wget https://github.com/akoova/MageComm/releases/download/v0.1.7-akoova/magecomm-darwin-arm64 && sudo mv magecomm-darwin-arm64 /usr/local/bin/magecomm && sudo chmod +x /usr/local/bin/magecomm`

If you wish to run a command for a clients environment use the following command: 
`MAGECOMM_ENVIRONMENT=ce23 aws-vault exec sudo-elasteradev -- magecomm magerun cache:clean`

## Release Tagging
As this is an akoova specific fork, we will be using akoova specific tag nomenclature with a semver postfix, e.g. `v1.0.0-akoova`

## Beta
Currently this tool is in beta, and is not recommended for production use.
Tested commands: RMQ based magerun command publishing, message listening, cat all supported archive types


## Installation

Install with binary and create configuration file in `/etc/magecomm/`(unix) / `%APPDATA%\magecomm\`(windows) or fallback to environment variables.
config file can be in yaml or json format e.g `config.yml` or `config.json`, please ensure you either set your n98 alias to magerun or explicitly set the path in the config file.
**WARNING: environment variables are not secure and can be easily read/modified by *any* user**

Download the latest release from the [releases page](https://github.com/furan917/magecomm/releases) for your platform and extract to a directory in your PATH.

example config.yml:
```
magecomm_log_path: /var/log/magecomm.log
magecomm_log_level: warn
magecomm_max_operational_cpu_limit: 80
magecomm_max_operational_memory_limit: 80
magecomm_environment: dev
magecomm_magerun_command_path: /usr/local/bin/n98-magerun2 --root-dir=/var/www/html
magecomm_listener_engine: sqs
magecomm_listener_allowed_queues:
  - magerun
magecomm_sqs_aws_region: eu-west-1
magecomm_rmq_tls: false
magecomm_rmq_user: guest
magecomm_rmq_pass: guest
magecomm_rmq_host: localhost
magecomm_rmq_port: 5672
magecomm_rmq_vhost: /
magecomm_slack_enabled: "true"
magecomm_slack_disable_output_notifications: "false"
magecomm_slack_webhook_url: https://hooks.slack.com/services/XXXXX/XXXXX/XXXXX
magecomm_slack_webhook_channel: "magecomm"
magecomm_slack_webhook_username: "magecomm"
magecomm_listeners:
  - magerun
magecomm_allowed_magerun_commands:
  - cache:clean
  - cache:flush
  - setup:static-content:deploy
  ...etc
magecomm_restricted_magerun_command_args:
  setup:static-content:deploy:
    - --jobs
  ...etc
magecomm_required_magerun_command_args:
    setup:upgrade:
        - --keep-generated
    ...etc
```

example config.json:
```
{
  "magecomm_log_path": "/var/log/magecomm.log",
  "magecomm_log_level": "warn",
  "magecomm_max_operational_cpu_limit": 80,
  "magecomm_max_operational_memory_limit": 80,
  "magecomm_environment": "dev",
  "magecomm_magerun_command_path": "/usr/local/bin/n98-magerun2 --root-dir=/var/www/html",
  "magecomm_listener_engine": "sqs",
  "magecomm_listener_allowed_queues": [
    "magerun"
  ],
  "magecomm_sqs_aws_region": "eu-west-1",
  "magecomm_rmq_tls": false,
  "magecomm_rmq_user": "guest",
  "magecomm_rmq_pass": "guest",
  "magecomm_rmq_host": "localhost",
  "magecomm_rmq_port": 5672,
  "magecomm_rmq_vhost": "/",
  "magecomm_slack_enabled": "true",
  "magecomm_slack_disable_output_notifications": "false",
  "magecomm_slack_webhook_url": "https://hooks.slack.com/services/XXXXX/XXXXX/XXXXX",
  "magecomm_slack_webhook_channel": "magecomm",
  "magecomm_slack_webhook_username": "magecomm",
  "magecomm_listeners": [
    "magerun"
  ],
  "magecomm_allowed_magerun_commands": [
    "cache:clean",
    "cache:flush",
    "setup:static-content:deploy"
    ...etc
  ],
  "magecomm_restricted_magerun_command_args": {
    "setup:static-content:deploy": [
      "--jobs"
    ]
    ...etc
  },
  "magecomm_required_magerun_command_args": {
    "setup:upgrade": [
      "--keep-generated"
    ]
    ...etc
  }
  ...etc
}
```

## Usage

### Global Flags

- `--debug`: Enable debug mode

e.g  
`magecomm --debug listen`  
`magecomm --debug magerun cache:clean`  
`magecomm --debug cat path/to/archive.tar.gz /path/to/file.txt`

### Commands

#### `magecomm magerun`

- A proxy for the magerun command via rmq/sqs with restricted command usage, allowed commands via `MAGECOMM_ALLOWED_MAGERUN_COMMANDS` with deeper control of args offered by MAGECOMM_RESTRICTED_MAGERUN_COMMAND_ARGS and MAGECOMM_REQUIRED_MAGERUN_COMMAND_ARGS
- Engine (sqs|rmq), default sqs, configured in config or by ENV `MAGECOMM_LISTENER_ENGINE`  
- The command will publish a message and then listen for the outputs return

#### `magecomm listen [queue1] [queue2] ...`

- Listen for messages from specified queues then handle them appropriately, fallback to config, then ENV `LISTENERS`
- Engine (sqs|rmq), default sqs, configured in config or by ENV `MAGECOMM_LISTENER_ENGINE`

#### `magecomm cat [archive] [filepath]`

- Extract a file from an archive and print its contents to stdout, we read headers to avoid being tricked by incorrect file extensions


## Configuration

The tool can be configured using a yaml or json config file at `/etc/magecomm/`(unix) / `%APPDATA%\magecomm\`(windows)  or by environment variables.
lowercase for file based config, uppercase for ENV.

The tool supports slack command run notifications via Webhook or App integration

- `MAGECOMM_LOG_PATH`: Path to log file, default: SYSLOG
- `MAGECOMM_LOG_LEVEL`: Log level, default: WARN, options (TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC)
- `MAGECOMM_MAX_OPERATIONAL_CPU_LIMIT`: Maximum CPU limit of system before we defer processing messages, default: 80
- `MAGECOMM_MAX_OPERATIONAL_MEMORY_LIMIT`: Maximum memory limit of system before we defer processing messages, default: 80
- `MAGECOMM_ENVIRONMENT`: the environment scope the tool is to work in, Default `default`
- `MAGECOMM_LISTENERS`: Comma-separated list of queues to listen to
- `MAGECOMM_LISTENER_ENGINE`: Listener engine to use (sqs/rmq), default: sqs
- `MAGECOMM_LISTENER_ALLOWED_QUEUES`: Comma-separated list of queues to allow to listen to, default: `cat, magerun`
- `MAGECOMM_PUBLISHER_OUTPUT_TIMEOUT`: Timeout for when listening to publisher message output return, default: 60s
- `MAGECOMM_MAGERUN_COMMAND_PATH` : Path to magerun command, default: `magerun` (expected alias of n98-magerun2.phar or /usr/local/bin/n98-magerun2 --root-dir=/magento/root/path) 
- `MAGECOMM_ALLOWED_MAGERUN_COMMANDS ` comma separated list of commands allowed to be run, fallback to in-code list
- `MAGECOMM_RESTRICTED_MAGERUN_COMMAND_ARGS` JSON object of commands and their restricted args, default: `{}`
- `MAGECOMM_REQUIRED_MAGERUN_COMMAND_ARGS` JSON object of commands and their required args, default: `{}`
- `MAGECOMM_SQS_AWS_REGION`: AWS region to use for SQS, default: eu-west-1
- `MAGECOMM_DEPLOY_ARCHIVE_PATH` path to the folder that contains the archives which are deployed, default: `/srv/magecomm/deploy/`
- `MAGECOMM_DEPLOY_ARCHIVE_LATEST_FILE` Filename of the latest archive (symlink), default: `latest.tar.gz`, if no value is set then MageComm will pick the latest created archive
- `MAGECOMM_RMQ_HOST` Default: `localhost`
- `MAGECOMM_RMQ_PORT` Default: `5672`
- `MAGECOMM_RMQ_USER` Default: ``
- `MAGECOMM_RMQ_PASS` Default: ``
- `MAGECOMM_RMQ_TLS`  Default: `false`
- `MAGECOMM_RMQ_VHOST` Default: `/`
- `MAGECOMM_SLACK_ENABLED` Default: `false`, (true|false), if true you must configure the WEBHOOK or APP configurations to work
- `MAGECOMM_SLACK_DISABLE_OUTPUT_NOTIFICATIONS` Default: 'false, (true|false), if true we will not send output notifications to slack
- `MAGECOMM_SLACK_WEBHOOK_URL` Default: ``
- `MAGECOMM_SLACK_WEBHOOK_CHANNEL` Default: ``
- `MAGECOMM_SLACK_WEBHOOK_USERNAME` Default: ``
- `MAGECOMM_SLACK_APP_TOKEN` Default: ``
- `MAGECOMM_SLACK_CHANNEL` Default: ``
- `MAGECOMM_SLACK_USERNAME` Default: ``

- ``

If using SQS the Pod/Instance this is deployed on must have an IAM role with the following permissions:
- `sqs:ReceiveMessage`
- `sqs:DeleteMessage`
- `sqs:DeleteMessageBatch`
- `sqs:GetQueueUrl`
- `sqs:ListQueues`
- `sqs:SendMessage`
- `sqs:SendMessageBatch`
- `sqs:GetQueueAttributes`
- `sqs:ChangeMessageVisibility`
- `sqs:ChangeMessageVisibilityBatch`
- `sqs:purgeQueue`
- `sts:GetCallerIdentity`

## Examples

1. Listen to messages from `magerun`:  
`magecomm listen magerun`  


2. Execute a magerun command using SQS as the publisher engine:  
`magecomm magerun cache:clean`  
`magecomm magerun setup:upgrade --keep-generated`


3. Extract and print the contents of a file from an archive (RAR, 7zip, and Xz are supported if installed)):  
`magecomm cat path/to/archive.zip /path/to/file.txt`  
`magecomm cat path/to/archive.tar /path/to/file.txt`  
`magecomm cat path/to/archive.tar.gz /path/to/file.txt`  
`magecomm cat path/to/archive.tar.bz2 /path/to/file.txt`  
`magecomm cat path/to/archive.tar.xz /path/to/file.txt`  
`magecomm cat path/to/archive.rar /path/to/file.txt`  
`magecomm cat path/to/archive.7zip /path/to/file.txt`
