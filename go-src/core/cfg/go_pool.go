package cfg

import (
	"encoding/json"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	GoPoolStr = "bufferSize.goPool"
)

type GoPoolCfg struct {
	Size int
}

var GoPoolCfgModel = &GoPoolCfg{}

func initGoPoolCfg() {
	data, _ := g.Cfg().GetWithCmd(gctx.New(), GoPoolStr)
	err := data.Scan(&GoPoolCfgModel)
	if err != nil {
		g.Log().Error(gctx.New(), "无法加载到协程池池大小配置数据")
	} else {
		cfgJson, _ := json.MarshalIndent(GoPoolCfgModel, "", " ")
		g.Log().Warningf(gctx.New(), "成功加载到协程池池大小配置数据:%s", cfgJson)
	}
}
