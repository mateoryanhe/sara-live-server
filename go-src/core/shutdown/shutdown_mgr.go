package shutdown

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
	"os"
	"time"
	"xr-game-server/constants/common"
	"xr-game-server/core/cfg"
)

const (
	ExitTime = 5
)

var commonHandler = make([]gproc.SigHandler, 0)
var loopHandler = make([]gproc.SigHandler, 0)

// ListenShutdown 监听程序退出
func ListenShutdown() {
	gproc.AddSigHandlerShutdown(doWhenAppExit)
	gproc.Listen()
}

// 处理程序退出
func doWhenAppExit(sig os.Signal) {
	exitTime := cfg.GetServerCfg()
	if exitTime.ExitTime == common.Zero {
		exitTime.ExitTime = ExitTime
	}
	//普通退出,占用时间少
	for _, handler := range commonHandler {
		handler(sig)
	}
	//数据库同步任务,防止程序退出,有数据没有同步到
	select {
	case <-time.After(time.Duration(exitTime.ExitTime) * time.Second):
		{
			g.Log().Warningf(gctx.New(), "程序正常退出,佛祖保佑")
		}
	}
}

// RegCommonShutDownHandler  注册普通程序退出处理器
func RegCommonShutDownHandler(handler ...gproc.SigHandler) {
	commonHandler = append(commonHandler, handler...)
}
