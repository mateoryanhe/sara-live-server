package cfg

func InitCfg() {
	initDefaultLoggerFromDetail()
	disableHTTPServerBuiltinLogs()
	initServerCfg()
	initStaticSiteCfg()
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
