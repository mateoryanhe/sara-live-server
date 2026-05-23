package recharge

// Init 服务启动时预加载充值配置缓存
func Init() {
	loadRechargeCfgCache()
	initGooglePay()
}
