package migrate

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	dsn := gstr.Replace(cfg.DefaultDbCfg.Link, "mysql:", "")
	g.Log().Warningf(gctx.New(), "同步表结构连接%s", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("表结构无法同步,连不上数据库")
	}
	// 迁移 schema
	for _, m := range dst {
		err = db.AutoMigrate(m)
	}
}
