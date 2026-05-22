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
	Level                uint32 `gorm:"uniqueIndex;default:0;comment:VIP等级" json:"level"`
	LevelName            string `gorm:"size:64;default:'';comment:等级名称" json:"levelName"`
	Status               uint8  `gorm:"default:0;comment:状态(0关闭,1开启)" json:"status"`
	UpgradeRechargeLimit uint64 `gorm:"default:0;comment:升级充值上限(美分,USD)" json:"upgradeRechargeLimit"`
	MinWithdrawAmount    uint64 `gorm:"default:0;comment:最低提现金额(美分,USD)" json:"minWithdrawAmount"`
	MaxWithdrawAmount    uint64 `gorm:"default:0;comment:最高提现金额(美分,USD)" json:"maxWithdrawAmount"`
	Fee                  uint64 `gorm:"default:0;comment:手续费(万分比,100=1%)" json:"fee"`
}

func initVipCfg() {
	migrate.AutoMigrate(&VipCfg{})
}
