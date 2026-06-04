package cfg

import (
	"github.com/gogf/gf/v2/frame/g"
)

func InitCfg() {
	//日志异步输出
	g.Log().SetAsync(true)
	initServerCfg()
	initDomainSiteCfg()
	initDbBufferCfg()
	initDbCfg()
	initGoPoolCfg()
	initWebSocketBufferCfg()
	initSensitiveWordCfg()
}
