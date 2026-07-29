package cfg

func InitCfg() {
	initDefaultLoggerFromDetail()
	initServerCfg()
	initDomainSiteCfg()
	initStaticPathCfg()
	initDbBufferCfg()
	initDbCfg()
	initDbQueryOnlyLogger()
	initGoPoolCfg()
	initWebSocketBufferCfg()
	initSensitiveWordCfg()
	logConfigContent()
}
