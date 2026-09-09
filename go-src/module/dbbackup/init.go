package dbbackup

import "xr-game-server/dao/cfgdao"

// Init 加载备份配置并订阅每日0点任务
func Init() {
	cfgdao.InitDbBackupCfgDao()
	cfgdao.ReloadDbBackupCfgCache()
	subscribeDayEvent()
}
