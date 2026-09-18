package entity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbStaticCacheRule db.TbName = "static_cache_rules"
)

// StaticCacheRule 静态网页不缓存文件名配置(CMS 管理,按文件名忽略大小写匹配).
type StaticCacheRule struct {
	migrate.OneModel
	FileName string `gorm:"size:255;uniqueIndex;default:'';comment:不缓存文件名(忽略大小写,不含目录)" json:"fileName"`
	Remark   string `gorm:"size:255;default:'';comment:备注" json:"remark"`
}

func initStaticCacheRule() {
	migrate.AutoMigrate(&StaticCacheRule{})
}
