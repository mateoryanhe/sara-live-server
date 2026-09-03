package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbAccountCfg db.TbName = "account_cfgs"
)

// 环境类型(账号配置)
const (
	AccountEnvTypeProd   uint8 = 0 // 正式服
	AccountEnvTypeReview uint8 = 1 // 提审服
	AccountEnvTypeTest   uint8 = 2 // 测试服
)

// AccountCfg 账号相关配置(CMS 管理,通常仅一条)
type AccountCfg struct {
	migrate.OneModel
	CancelAccountByCodeEnabled bool  `gorm:"default:0;comment:注销码销户开关(官网公开接口)" json:"cancelAccountByCodeEnabled"`
	BlockSimulatorLogin        bool  `gorm:"default:0;comment:拦截模拟器登录(默认关闭=不拦截)" json:"blockSimulatorLogin"`
	EnvType                    uint8 `gorm:"default:0;comment:环境类型(0正式服,1提审服,2测试服)" json:"envType"`
}

func initAccountCfg() {
	migrate.AutoMigrate(&AccountCfg{})
}
