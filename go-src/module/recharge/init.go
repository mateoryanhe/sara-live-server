package recharge

import "xr-game-server/dao/cfgdao"

// Init 服务启动时预加载充值配置缓存
func Init() {
	loadRechargeCfgCache()
	initGooglePlayCfg()
	cfgdao.InitHaiPayCollectionCountryCfgDao()
	cfgdao.ReloadHaiPayCollectionCountryCfgCache()
	cfgdao.InitHaiPayCoinMerchantCollectionCfgDao()
	cfgdao.ReloadHaiPayCoinMerchantCollectionCfgCache()
	initHaiPayCfg()
	initRechargeOrderTimeoutWatch()
}
