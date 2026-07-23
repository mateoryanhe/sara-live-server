package cfg

func InitCfg() {
	initDefaultLoggerFromDetail()
	initServerCfg()
	initDomainSiteCfg()
	initStaticPathCfg()
	initDbBufferCfg()
	initDbCfg()
	initGoPoolCfg()
	initWebSocketBufferCfg()
	initSensitiveWordCfg()
	logConfigContent()
}
