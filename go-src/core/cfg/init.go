package cfg

func InitCfg() {
	initDefaultLoggerFromDetail()
	disableHTTPServerBuiltinLogs()
	initServerCfg()
	initDomainSiteCfg()
	initStaticPathCfg()
	initCMSFileExportCfg()
	initDbBufferCfg()
	initDbCfg()
	initDbQueryOnlyLogger()
	initGoPoolCfg()
	initWebSocketBufferCfg()
	initSensitiveWordCfg()
	initIpGeoCfg()
	logConfigContent()
}
