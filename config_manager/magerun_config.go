package config_manager

import (
	"magecomm/logger"
	"strings"
)

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

func IsMageRunCommandAllowed(command string) bool {
	var allowedCommands []string

	allowedCommandsConfig := GetValue(CommandConfigAllowedMageRunCommands)
	if allowedCommandsConfig != "" {
		allowedCommands = strings.Split(allowedCommandsConfig, ",")
	} else {
		allowedCommands = defaultAllowedCommands
	}

	for _, allowedCommand := range allowedCommands {
		if allowedCommand == command {
			return true
		}
	}
	// print allowed commands
	logger.Errorf("`%s` Command not allowed, allowed commands are:\n%s \n", command, strings.Join(allowedCommands, ",\n"))
	return false
}

func IsRestrictedCommandArgsIncluded(command string, args []string) bool {
	restrictedArgsString := GetValue(CommandConfigRestrictedMagerunCommandArgs)
	if restrictedArgsString == "" {
		return false
	}
	restrictedCommandArgMap := ParseCommandArgsMap(restrictedArgsString)

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
	requiredArgsString := GetValue(CommandConfigRequiredMagerunCommandArgs)
	if requiredArgsString == "" {
		return true, []string{}
	}
	requiredCommandArgMap := ParseCommandArgsMap(requiredArgsString)
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
