package entity

// Init 充值相关表迁移与 syndb 注册
func Init() {
	initRechargeCfg()
	initRechargeOrder()
	initGooglePlayCfg()
}
