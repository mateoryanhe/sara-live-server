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
	AppId              string `gorm:"size:64;default:'';comment:声网AppId" json:"appId"`
	AppCertificate     string `gorm:"size:128;default:'';comment:声网AppCertificate" json:"appCertificate"`
	TokenExpireSeconds uint32 `gorm:"default:21600;comment:Token有效期(秒,默认6小时)" json:"tokenExpireSeconds"`
}

func initAgoraCfg() {
	migrate.AutoMigrate(&AgoraCfg{})
}
