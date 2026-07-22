package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbAppPkg db.TbName = "app_pkgs"
)

// AppPkg App包配置(CMS管理)
type AppPkg struct {
	migrate.OneModel
	PackageName string `gorm:"uniqueIndex;size:128;comment:包名" json:"packageName"`
	SecretKey   string `gorm:"size:256;comment:密钥" json:"secretKey"`
	Remark      string `gorm:"size:512;default:'';comment:备注" json:"remark"`
}

func initAppPkg() {
	migrate.AutoMigrate(&AppPkg{})
}
