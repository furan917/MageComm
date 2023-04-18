# MageComm CLI Tool

MageComm CLI is a command line tool for managing Magento applications and deployments. It provides a convenient way to execute restricted magerun commands, manage deployments, and cat specific files within an archive.
The main use case for this tool is to provide a way to execute magerun commands in a controlled manner via a messaging service. This allows us to execute commands on a remote server without exposing the application server itself

*It is important to note that the environment configuration/env is important to set if you plan to use this tool in a shared rmq/sqs instance as this will prefix your queues to avoid cross communication*

The tool looks for a configuration file in `/etc/magecomm/`(unix) / `%APPDATA%\magecomm\`(windows) `config.yml|config.json` and if none found defaults to environment variables, then fallback to default values.

## Beta
Currently this tool is in beta, and is not recommended for production use.
Tested commands: RMQ based magerun command publishing, message listening, cat all supported archive types


## Installation

Install with binary and create configuration file in `/etc/magecomm/`(unix) / `%APPDATA%\magecomm\`(windows) or fallback to environment variables.
config file can be in yaml or json format e.g `config.yml` or `config.json`
**WARNING: environment variables are not secure and can be easily read/modified by *any* user**

Download the latest release from the [releases page](https://github.com/furan917/magecomm/releases) for your platform and extract to a directory in your PATH.

example config.yml:
```
magecomm_log_path: /var/log/magecomm.log
magecomm_log_level: warn
magecomm_max_operational_cpu_limit: 80
magecomm_max_operational_memory_limit: 80
magecomm_environment: dev
magecomm_listener_engine: sqs
sqs_aws_region: eu-west-1
rmq_tls: false
rmq_user: guest
rmq_pass: guest
rmq_host: localhost
rmq_port: 5672
rmq_vhost: /
magecomm_listeners:
  - magerun
  - deploy
magecomm_allowed_magerun_commands:
  - cache:clean
  - cache:flush
  - setup:static-content:deploy
  ...etc
```

example config.json:
```
{
  "magecomm_log_path": /var/log/magecomm.log
  "magecomm_log_level": warn
  "magecomm_max_operational_cpu_limit": 80,
  "magecomm_max_operational_memory_limit": 80,
  "magecomm_environment": "dev",
  "magecomm_listener_engine": "sqs",
  "sqs_aws_region": "eu-west-1",
  "rmq_tls": false,
  "rmq_user": "guest",
  "rmq_pass": "guest",
  "rmq_host": "localhost",
  "rmq_port": 5672,
  "rmq_vhost": "/",
  "magecomm_listeners": [
    "magerun",
    "deploy"
  ],
  "magecomm_allowed_magerun_commands": [
    "cache:clean",
    "cache:flush",
    "setup:static-content:deploy"
  ],
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

- A proxy for the magerun command via rmq/sqs with restricted command usage, allowed commands via `MAGECOMM_ALLOWED_MAGERUN_COMMANDS`
- Engine (sqs|rmq), default sqs, configured in config or by ENV `MAGECOMM_LISTENER_ENGINE`  
- The command will publish a message and then listen for the outputs return

#### `magecomm listen [queue1] [queue2] ...`

- Listen for messages from specified queues then handle them appropriately, fallback to config, then ENV `LISTENERS`
- Engine (sqs|rmq), default sqs, configured in config or by ENV `MAGECOMM_LISTENER_ENGINE`

#### `magecomm deploy [filepath]`

- (WIP) Deploy an archived file to the specified environment

#### `magecomm cat-deploy [filepath]`

- Extract a file from the latest deployed archive and print its contents to stdout.  
  *Command must have at minimum the `MAGECOMM_DEPLOY_ARCHIVE_PATH` configured in config file or by ENV to work*

#### `magecomm cat [archive] [filepath]`

- Extract a file from an archive and print its contents to stdout, we read headers to avoid being tricked by incorrect file extensions


## Configuration

The tool can be configured using a yaml or json config file at `/etc/magecomm/`(unix) / `%APPDATA%\magecomm\`(windows)  or by environment variables.
lowercase for file based config, uppercase for ENV

- `MAGECOMM_LOG_PATH`: Path to log file, default: SYSLOG
- `MAGECOMM_LOG_LEVEL`: Log level, default: WARN, options (TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC)
- `MAGECOMM_MAX_OPERATIONAL_CPU_LIMIT`: Maximum CPU limit of system before we defer processing messages, default: 80
- `MAGECOMM_MAX_OPERATIONAL_MEMORY_LIMIT`: Maximum memory limit of system before we defer processing messages, default: 80
- `MAGECOMM_ENVIRONMENT`: the environment scope the tool is to work in, Default `default`
- `MAGECOMM_LISTENERS`: Comma-separated list of queues to listen to
- `MAGECOMM_LISTENER_ENGINE`: Listener engine to use (sqs/rmq), default: sqs
- `MAGECOMM_PUBLISHER_OUTPUT_TIMEOUT`: Timeout for when listening to publisher message output return, default: 60s
- `MAGECOMM_ALLOWED_MAGERUN_COMMANDS ` comma separated list of commands allowed to be run, fallback to in-code list
- `SQS_AWS_REGION`: AWS region to use for SQS, default: eu-west-1
- `DEPLOY_ARCHIVE_PATH` path to the folder that contains the archives which are deployed, default: `/srv/magecomm/deploy/`
- `DEPLOY_ARCHIVE_LATEST_FILE` Filename of the latest archive (symlink), default: `latest.tar.gz`, if no value is set then MageComm will pick the latest created archive
- `RMQ_HOST` Default: `localhost`
- `RMQ_PORT` Default: `5672`
- `RMQ_USER` Default: ``
- `RMQ_PASS` Default: ``
- `RMQ_TLS`  Default: `false`
- `RMQ_VHOST` Default: `/`

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

1. Listen to messages from `magerun` and `deploy`:  
`magecomm listen magerun deploy`  


2. Execute a magerun command using SQS as the publisher engine:  
`magecomm magerun cache:clean`  
`magecomm magerun setup:upgrade --keep-generated`  


3. Deploy a gzipped file:  
`magecomm deploy path/to/archive.gz`    


4. Extract and print the contents of a file from an archive (RAR, 7zip, and Xz are supported if installed)):  
`magecomm cat path/to/archive.zip /path/to/file.txt`  
`magecomm cat path/to/archive.tar /path/to/file.txt`  
`magecomm cat path/to/archive.tar.gz /path/to/file.txt`  
`magecomm cat path/to/archive.tar.bz2 /path/to/file.txt`  
`magecomm cat path/to/archive.tar.xz /path/to/file.txt`  
`magecomm cat path/to/archive.rar /path/to/file.txt`  
`magecomm cat path/to/archive.7zip /path/to/file.txt`  


5. Extract and print the contents of a file from the latest deploy  
`magecomm cat-deploy /path/to/target/file.txt`  