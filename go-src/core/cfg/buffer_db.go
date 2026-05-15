package cfg

import (
	"encoding/json"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	LazyDbStr = "bufferSize.db.lazy"
	FastDbStr = "bufferSize.db.fast"
)

type DbBufferCfg struct {
	//单行数据配置
	Period int
	Small  int
	Middle int
	Large  int
}

var LazyDbBufferCfg = &DbBufferCfg{}
var FastDbBufferCfg = &DbBufferCfg{}

func initDbBufferCfg() {
	initLazy()
	initFast()
}

func initLazy() {
	data, _ := g.Cfg().GetWithCmd(gctx.New(), LazyDbStr)
	err := data.Scan(&LazyDbBufferCfg)
	if err != nil {
		g.Log().Error(gctx.New(), "无法加载到延迟缓冲池大小配置数据")
	} else {
		cfgJson, _ := json.MarshalIndent(LazyDbBufferCfg, "", " ")
		g.Log().Warningf(gctx.New(), "成功加载到延迟缓冲池大小配置数据:%s", cfgJson)
	}
}

func initFast() {
	data, _ := g.Cfg().GetWithCmd(gctx.New(), FastDbStr)
	err := data.Scan(&FastDbBufferCfg)
	if err != nil {
		g.Log().Error(gctx.New(), "无法加载到快速缓冲池大小配置数据")
	} else {
		cfgJson, _ := json.MarshalIndent(FastDbBufferCfg, "", " ")
		g.Log().Warningf(gctx.New(), "成功加载到快速缓冲池大小配置数据:%s", cfgJson)
	}
}
