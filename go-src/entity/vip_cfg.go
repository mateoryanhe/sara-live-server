package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbVipCfg db.TbName = "vip_cfgs"
)

const (
	VipCfgStatusDisabled uint8 = 0 // 关闭
	VipCfgStatusEnabled  uint8 = 1 // 开启
)

// VipCfg VIP等级配置(CMS管理)
type VipCfg struct {
	migrate.OneModel
	Level                uint32  `gorm:"uniqueIndex;default:0;comment:VIP等级" json:"level"`
	LevelName            string  `gorm:"size:64;default:'';comment:等级名称" json:"levelName"`
	Status               uint8   `gorm:"default:0;comment:状态(0关闭,1开启)" json:"status"`
	UpgradeRechargeLimit float64 `gorm:"type:decimal(18,4);default:0;comment:升级充值上限(USD)" json:"upgradeRechargeLimit"`
	MinWithdrawAmount    float64 `gorm:"type:decimal(18,4);default:0;comment:最低提现金额(USD)" json:"minWithdrawAmount"`
	MaxWithdrawAmount    float64 `gorm:"type:decimal(18,4);default:0;comment:最高提现金额(USD)" json:"maxWithdrawAmount"`
	Fee                  float64 `gorm:"type:decimal(18,4);default:0;comment:手续费" json:"fee"`
}

func initVipCfg() {
	migrate.AutoMigrate(&VipCfg{})
}
