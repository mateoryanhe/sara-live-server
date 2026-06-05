package cfg

func InitCfg() {
	//日志异步输出
	initServerCfg()
	initDomainSiteCfg()
	initDbBufferCfg()
	initDbCfg()
	initGoPoolCfg()
	initWebSocketBufferCfg()
	initSensitiveWordCfg()
}
