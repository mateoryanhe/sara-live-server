package startup

import (
	"context"
	"time"

	"xr-game-server/core/xrlog"
)

var startedAt time.Time

// MarkBegin 记录进程启动起点(在 main 入口尽早调用).
func MarkBegin() {
	startedAt = time.Now()
}

// LogEnd 输出启动成功及耗时(HTTP 服务就绪后调用).
func LogEnd(ctx context.Context) {
	ms := time.Since(startedAt).Milliseconds()
	xrlog.DetailLog.Infof(ctx, "佛祖保佑，启动成功，耗时 %dms", ms)
}
