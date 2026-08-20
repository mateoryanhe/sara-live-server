package entity

// Init 游戏相关表迁移与 syndb 注册
func Init() {
	initGameCfg()
	initVendorGameLib()
	initGameBetLog()
	initGameWinLog()
	initGamePlatformCfg()
}
