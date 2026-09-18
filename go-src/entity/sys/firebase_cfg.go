package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbFirebaseCfg db.TbName = "firebase_cfgs"
)

// FirebaseCfg Firebase Authentication 服务端验签配置(CMS 管理,通常仅一条).
type FirebaseCfg struct {
	migrate.OneModel
	ProjectId          string `gorm:"size:128;default:'';comment:Firebase Project ID" json:"projectId"`
	ClientConfigJson   string `gorm:"type:text;comment:Firebase客户端配置JSON" json:"clientConfigJson"`
	ServiceAccountJson string `gorm:"type:text;comment:Firebase服务账号JSON" json:"serviceAccountJson"`
}

func (FirebaseCfg) TableName() string {
	return string(TbFirebaseCfg)
}

func initFirebaseCfg() {
	migrate.AutoMigrate(&FirebaseCfg{})
}
