package cfg

func InitCfg() {
	initDefaultLoggerFromDetail()
	initServerCfg()
	initDomainSiteCfg()
	initDbBufferCfg()
	initDbCfg()
	initGoPoolCfg()
	initWebSocketBufferCfg()
	initSensitiveWordCfg()
	logConfigContent()
}
