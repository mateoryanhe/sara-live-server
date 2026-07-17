package cfg

func InitCfg() {
	initServerCfg()
	initDomainSiteCfg()
	initDbBufferCfg()
	initDbCfg()
	initGoPoolCfg()
	initWebSocketBufferCfg()
	initSensitiveWordCfg()
	logConfigContent()
}
