package cfg

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	DefaultDbCfgStr = "database.default"
)

type DBCfg struct {
	Link        string
	Extra       string
	Debug       string
	MaxIdle     int           `json:"maxIdle"`
	MaxOpen     int           `json:"maxOpen"`
	MaxLifetime time.Duration `json:"maxLifetime"`
}

var DefaultDbCfg = &DBCfg{}

func initDbCfg() {
	data, _ := g.Cfg().GetWithCmd(gctx.New(), DefaultDbCfgStr)
	err := data.Scan(&DefaultDbCfg)
	if err != nil {
		g.Log().Error(gctx.New(), "无法加载到数据库配置数据")
	}
}
