package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbAccountCfg db.TbName = "account_cfgs"
)

// AccountCfg 账号相关配置(CMS 管理,通常仅一条)
type AccountCfg struct {
	migrate.OneModel
	CancelAccountByCodeEnabled bool `gorm:"default:0;comment:注销码销户开关(官网公开接口)" json:"cancelAccountByCodeEnabled"`
	SimulatorLoginEnabled      bool `gorm:"default:1;comment:模拟器是否允许登录" json:"simulatorLoginEnabled"`
}

func initAccountCfg() {
	migrate.AutoMigrate(&AccountCfg{})
}
