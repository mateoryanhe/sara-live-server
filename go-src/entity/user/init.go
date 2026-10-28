package entity

// Init 用户/账号相关表迁移与 syndb 注册
func Init() {
	initAccount()
	initUserInfo()
	initUserLoginDevice()
	initUserLoginLocation()
	initUserExt()
	initUserRechargeCfgFirstRecharge()
	initUserCumulativeStat()
	initCurrencyLog()
	initCoinMerchantGoldTransferLog()
	initUserAttributionReport()
	initUserFeedbackRecord()
	initUserReportRecord()
	initRandomNickname()
	initWalletExchangeCfg()
	initCustomerServiceCfg()
	initPrivacyPolicyCfg()
	initSofiePrivacyPolicyCfg()
	initShareAppDownloadCfg()
	initAppToken()
	initUserMaxId()
	initUserPayoutInfo()
	initUserWithdrawLog()
	initUserWithdrawTierCfg()
}
