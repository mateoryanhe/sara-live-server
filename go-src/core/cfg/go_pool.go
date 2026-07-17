package cfg

import (
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
	}
}
