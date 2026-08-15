package entity

// Init 用户/账号相关表迁移与 syndb 注册
func Init() {
	initAccount()
	initUserInfo()
	initUserLoginDevice()
	initUserExt()
	initUserCumulativeStat()
	initCurrencyLog()
	initRandomNickname()
	initWalletExchangeCfg()
	initCustomerServiceCfg()
	initPrivacyPolicyCfg()
	initAppToken()
}
