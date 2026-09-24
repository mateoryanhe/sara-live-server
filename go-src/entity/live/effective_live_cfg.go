package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const TbEffectiveLiveCfg db.TbName = "effective_live_cfgs"

// EffectiveLiveCfg 有效直播时长配置(CMS 管理,通常仅一条记录).
type EffectiveLiveCfg struct {
	migrate.OneModel
	MinSessionMinutes int `gorm:"default:30;comment:单场直播计入有效时长的门槛(分钟)" json:"minSessionMinutes"`
}

func initEffectiveLiveCfg() {
	migrate.AutoMigrate(&EffectiveLiveCfg{})
}
