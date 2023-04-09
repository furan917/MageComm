package system_limits

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"magecomm/config_manager"
	"magecomm/logger"
	"strconv"
	"time"
)

const WaitTimeBetweenChecks = 15 * time.Second

func CheckIfOutsideOperationalLimits() bool {
	return isCpuLimitReached() || isMemoryLimitReached()
}

func isCpuLimitReached() bool {
	cpuLimitStr := config_manager.GetValue(config_manager.CommandConfigMaxOperationalCpuLimit)
	cpuLimit, err := strconv.ParseFloat(cpuLimitStr, 64)
	if err != nil {
		return false
	}

	percent, err := cpu.Percent(0, false)
	if err != nil {
		return false
	}

	if len(percent) > 0 && percent[0] > cpuLimit {
		logger.Warnf("Max CPU limit reached. Current CPU usage: %f", percent[0])
		return true
	}

	return false
}

func isMemoryLimitReached() bool {
	memoryLimitStr := config_manager.GetValue(config_manager.CommandConfigMaxOperationalMemoryLimit)
	memoryLimit, err := strconv.ParseUint(memoryLimitStr, 10, 64)
	if err != nil {
		return false
	}

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return false
	}

	usedMemory := vmStat.UsedPercent
	if usedMemory > float64(memoryLimit) {
		logger.Warnf("Max memory limit reached. Current memory usage: %f", usedMemory)
		return true
	}

	return false
}

func SystemLimitCheckSleep() {
	logger.Warnf("Outside operational limits. Waiting %s seconds then will check again...", WaitTimeBetweenChecks/time.Second)
	time.Sleep(WaitTimeBetweenChecks)
}
