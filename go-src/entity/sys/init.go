package entity

// Init 系统配置与资源监控相关表迁移与 syndb 注册
func Init() {
	initUploadResourceCfg()
	initPreloadCfg()
	initDataSyncCfg()
	initH5LiveDeployCfg()
	initSysResourceMetric()
	initSysResourceMetricAgg()
}
