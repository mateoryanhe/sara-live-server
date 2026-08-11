package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbDataSyncCfg db.TbName = "data_sync_cfgs"
)

// DataSyncCfg 数据同步配置(CMS 管理,通常仅一条;运行时直查数据库)
type DataSyncCfg struct {
	migrate.OneModel
	TargetApiBase string `gorm:"size:512;default:'';comment:目标API根地址" json:"targetApiBase"`
	Token         string `gorm:"size:512;default:'';comment:同步Token" json:"token"`
}

func initDataSyncCfg() {
	migrate.AutoMigrate(&DataSyncCfg{})
}
