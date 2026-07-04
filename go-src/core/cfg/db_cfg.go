package cfg

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/xrjson"
)

const (
	DefaultDbCfgStr = "database.default"
)

type DBCfg struct {
	Link  string
	Extra string
	Debug string
}

var DefaultDbCfg = &DBCfg{}

func initDbCfg() {
	data, _ := g.Cfg().GetWithCmd(gctx.New(), DefaultDbCfgStr)
	err := data.Scan(&DefaultDbCfg)
	if err != nil {
		g.Log().Error(gctx.New(), "无法加载到数据库配置数据")
	} else {
		cfgJson := xrjson.MustMarshalIndent(DefaultDbCfg)
		g.Log().Warningf(gctx.New(), "成功加载到数据库配置数据:%s", cfgJson)
	}
}
