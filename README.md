# MageComm CLI Tool

MageComm CLI is a command line tool for managing Magento applications and deployments. It provides a convenient way to execute restricted magerun commands, manage deployments, interact with message queues (RabbitMQ or SQS), and cat specific files within an archive.

*It is important to note that the environment configuration/env is important to set if you plan to use this tool in a shared rmq/sqs instance as this will prefix your queues to avoid cross communication*

The tool looks for a configuration file in `/etc/magecomm/config.yml`(unix) / `%APPDATA%\magecomm\config.yml`(windows) and if none found defaults to environment variables, then fallback to default values.

## Beta
Currently this tool is in beta, and is not recommended for production use.
Tested commands: RMQ based magerun command publishing, message listening, cat all supported archive types


## Installation

Install with binary and create configuration file in `/etc/magecomm/config.yml`(unix) / `%APPDATA%\magecomm\config.yml`(windows) or fallback to environment variables:  
**WARNING: environment variables are not secure and can be easily read/modified by *any* user**

Download the latest release from the [releases page](https://github.com/furan917/magecomm/releases) for your platform and extract to a directory in your PATH.

example config.yml:
```
environment: dev
listener_engine: sqs
sqs_aws_region: eu-west-1
rmq_tls: false
rmq_user: guest
rmq_pass: guest
rmq_host: localhost
rmq_port: 5672
rmq_vhost: /
listeners:
  - magerun
  - deploy
allowed_magerun_commands:
  - cache:clean
  - cache:flush
  - setup:static-content:deploy
  ...etc
```

## Usage

### Global Flags

- `--debug`: Enable debug mode

e.g  
`magecomm --debug listen`  
`magecomm --debug magerun cache:clean`
`magecomm --debug cat path/to/archive.tar.gz /path/to/file.txt`

### Commands

#### `magecomm listen [queue1] [queue2] ...`

- Listen for messages from specified queues then handle them appropriately, fallback to config, then ENV LISTENERS
- Engine (sqs|rmq), default sqs, configured in config or by ENV LISTENER_ENGINE

#### `magecomm magerun`

- A proxy for the magerun command via rmq/sqs with restricted command usage, fallback to config, then ENV LISTENERS
- Engine (sqs|rmq), default sqs, configured in config or by ENV LISTENER_ENGINE

#### `magecomm deploy [filepath]`

- (WIP) Deploy an archived file to the specified environment

#### `magecomm cat [archive] [filepath]`

- Extract a file from an archive and print its contents to stdout, we read headers to avoid being tricked by incorrect file extensions

## Configuration

The tool can be configured using a yaml file at `/etc/magecomm/config.yml`(unix) / `%APPDATA%\magecomm\config.yml`(windows)  or by environment variables.
lowercase for yml, uppercase for ENV

- `ENVIRONMENT`: the environment scope the tool is to work in, Default `default`
- `LISTENERS`: Comma-separated list of queues to listen to
- `LISTENER_ENGINE`: Listener engine to use (sqs/rmq), default: sqs
- `ALLOWED_MAGERUN_COMMANDS`: Comma-separated list of magerun commands to allow
- `SQS_AWS_REGION`: AWS region to use for SQS, default: eu-west-1
- `ALLOWED_MAGERUN_COMMANDS ` comma separated list of commands allowed to be run, fallback to in-code list
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

