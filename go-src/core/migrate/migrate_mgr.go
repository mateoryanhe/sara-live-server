package migrate

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"time"
	"xr-game-server/core/cfg"
)

type OneModel struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type MoreModel struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	IsDeleted bool      `gorm:"default:0" json:"-"`
	DeletedAt time.Time `json:"-"`
}

// AutoMigrate 代码启动,自动同步表结构,表结构保持驼峰
func AutoMigrate(dst ...interface{}) {
	link := cfg.DefaultDbCfg.Link
	if !strings.HasPrefix(link, "pgsql:") {
		panic("database.default.link must use pgsql driver")
	}
	dsn := PostgresDSN(link, cfg.DefaultDbCfg.Extra)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("表结构无法同步,连不上数据库")
	}
	// 迁移 schema
	for _, m := range dst {
		err = db.AutoMigrate(m)
	}
}
