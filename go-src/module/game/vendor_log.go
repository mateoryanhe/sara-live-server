package game

import (
	"github.com/gogf/gf/v2/os/glog"
	"xr-game-server/core/xrlog"
)

func vendorDetailLog() *glog.Logger {
	return xrlog.DetailLog.Logger()
}
