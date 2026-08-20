package hotrestart

import (
	"context"

	"xr-game-server/core/syndb"
	"xr-game-server/core/xrlog"
)

const (
	msgPhase1Start       = "热更新---第一阶段----开始"
	msgPhase1QueueEmpty  = "热更新---第一阶段---入库队列清空完毕"
	msgPhase1End         = "热更新----第一阶段-----结束"
	msgPhase2Start       = "热更新---第二阶段----开始"
	msgPhase2HTTPClose   = "热更新---第二阶段----收到http请求,继续执行,URI=%s ip=%s"
	msgPhase2QueueHasRow = "热更新---第二阶段----入库队列有数据,数据量%d行,%s"
	msgPhase2EndBless    = "热更新---第二阶段----结束,佛祖保佑"
	msgPhase2EndNoBless  = "热更新---第二阶段----结束,佛祖不保佑了,要去看看日志,%s"
)

func logPhase1Start(ctx context.Context) {
	xrlog.DetailLog.Warning(ctx, msgPhase1Start)
}

func logPhase1QueueEmpty(ctx context.Context) {
	xrlog.DetailLog.Warning(ctx, msgPhase1QueueEmpty)
}

func logPhase1End(ctx context.Context) {
	xrlog.DetailLog.Warning(ctx, msgPhase1End)
}

func logPhase2Start(ctx context.Context) {
	xrlog.DetailLog.Warning(ctx, msgPhase2Start)
}

// LogPhase2HTTPRequest 第二阶段旧进程收到 HTTP/WS 请求(继续执行,仅记录日志).
func LogPhase2HTTPRequest(ctx context.Context, uri, ip string) {
	xrlog.DetailLog.Warningf(ctx, msgPhase2HTTPClose, uri, ip)
}

func logPhase2End(ctx context.Context) {
	rowCount := syndb.PendingQueueRowCount()
	summary := syndb.PendingQueueSummary()
	if rowCount > 0 {
		if summary == "" {
			summary = "detail=unknown"
		}
		xrlog.DetailLog.Warningf(ctx, msgPhase2QueueHasRow, rowCount, summary)
		xrlog.DetailLog.Warningf(ctx, msgPhase2EndNoBless, summary)
		return
	}
	xrlog.DetailLog.Warning(ctx, msgPhase2EndBless)
}
