package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbVendorGameLib db.TbName = "vendor_game_libs"
)

// VendorGameLib 第三方游戏库(CMS 同步入库, 查询读此表).
type VendorGameLib struct {
	migrate.OneModel
	GameCode string `gorm:"uniqueIndex:idx_vendor_game_lib_code_platform;size:64;comment:游戏编码" json:"gameCode"`
	Name     string `gorm:"size:128;default:'';comment:中文名称" json:"name"`
	NameEn   string `gorm:"size:128;default:'';comment:英文名称" json:"nameEn"`
	Category string `gorm:"size:64;default:'';comment:分类" json:"category"`
	Cover    string `gorm:"size:512;default:'';comment:封面相对路径" json:"cover"`
	Platform string `gorm:"uniqueIndex:idx_vendor_game_lib_code_platform;size:32;default:'';comment:平台编码" json:"platform"`
}

func initVendorGameLib() {
	migrate.AutoMigrate(&VendorGameLib{})
	migrate.DropIndexIfExists(string(TbVendorGameLib), "idx_vendor_game_libs_game_code")
}
