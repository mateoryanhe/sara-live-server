package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbAgoraCfg db.TbName = "agora_cfgs"
)

// AgoraCfg 声网配置(CMS管理,通常仅一条记录)
type AgoraCfg struct {
	migrate.OneModel
	AppId               string `gorm:"size:64;default:'';comment:声网AppId" json:"appId"`
	AppCertificate      string `gorm:"size:128;default:'';comment:声网AppCertificate" json:"appCertificate"`
	RestCustomerId      string `gorm:"size:64;default:'';comment:声网REST CustomerId" json:"restCustomerId"`
	RestCustomerSecret  string `gorm:"size:128;default:'';comment:声网REST CustomerSecret" json:"restCustomerSecret"`
	CloudPlayerRegion   string `gorm:"size:16;default:'cn';comment:云播放器区域(cn/ap/eu/na)" json:"cloudPlayerRegion"`
	TokenExpireSeconds  uint32 `gorm:"default:86400;comment:声网Token有效期(秒,4-24小时,默认24小时)" json:"tokenExpireSeconds"`
	TokenRefreshSeconds uint32 `gorm:"default:79200;comment:Token提前刷新阈值(秒,2小时起,需比有效期至少少2小时)" json:"tokenRefreshSeconds"`
}

func initAgoraCfg() {
	migrate.AutoMigrate(&AgoraCfg{})
}
