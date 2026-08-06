package cfg

func InitCfg() {
	initDefaultLoggerFromDetail()
	disableHTTPServerBuiltinLogs()
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
