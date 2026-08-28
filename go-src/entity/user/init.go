package entity

// Init 用户/账号相关表迁移与 syndb 注册
func Init() {
	initAccount()
	initUserInfo()
	initUserLoginDevice()
	initUserExt()
	initUserRechargeCfgFirstRecharge()
	initUserCumulativeStat()
	initCurrencyLog()
	initRandomNickname()
	initWalletExchangeCfg()
	initCustomerServiceCfg()
	initPrivacyPolicyCfg()
	initAppToken()
	initUserMaxId()
}
