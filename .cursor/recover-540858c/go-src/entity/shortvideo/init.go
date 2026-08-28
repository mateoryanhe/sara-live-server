package entity

// Init 短视频相关表迁移与 syndb 注册
func Init() {
	initShortVideo()
	initShortVideoStat()
	initShortVideoAuthorStat()
	initShortVideoCfg()
	initShortVideoCategory()
	initShortVideoWatch()
	initShortVideoPriceTier()
	initShortVideoAuthorSettlementLog()
}
