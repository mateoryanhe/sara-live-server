package logquery

import (
	"sync"

	appcfg "xr-game-server/core/cfg"
)

var (
	logQueryConfigMu    sync.RWMutex
	logQueryConfigCache logQueryConfig
	logQueryConfigReady bool
)

// Init 启动时加载并缓存日志查询配置,避免定时任务中重复 MustGet 触发配置文件查找失败
func Init() {
	refreshLogQueryConfig()
}

func loadLogQueryConfig() logQueryConfig {
	logQueryConfigMu.RLock()
	if logQueryConfigReady {
		cfg := logQueryConfigCache
		logQueryConfigMu.RUnlock()
		return cfg
	}
	logQueryConfigMu.RUnlock()
	return refreshLogQueryConfig()
}

func refreshLogQueryConfig() logQueryConfig {
	cfg := buildLogQueryConfig().normalized()
	logQueryConfigMu.Lock()
	defer logQueryConfigMu.Unlock()
	if cfg.hasEssentialPaths() {
		logQueryConfigCache = cfg
		logQueryConfigReady = true
		return cfg
	}
	if logQueryConfigReady {
		return logQueryConfigCache
	}
	logQueryConfigCache = cfg
	return cfg
}

func (c logQueryConfig) hasEssentialPaths() bool {
	return c.LogDir != "" || c.exportAbsDir() != "" || appcfg.GetCMSFileExportRoot() != ""
}
