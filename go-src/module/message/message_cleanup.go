package message

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/dao/messagedao"
	"xr-game-server/gameevent"
)

func initMessageCleanup() {
	event.Sub(gameevent.DayEvent, onDayCleanupUserMessages)
}

func onDayCleanupUserMessages(_ any) {
	ctx := gctx.New()
	g.Log().Info(ctx, "开始执行 user_messages 超量清理任务")
	deletedMessages, deletedSessions := messagedao.CleanupExcessUserMessages(ctx)
	if deletedMessages == 0 && deletedSessions == 0 {
		g.Log().Info(ctx, "user_messages 未超过10万条,无需清理")
		return
	}
	g.Log().Infof(ctx, "user_messages 超量清理完成,删除消息=%d,删除session=%d", deletedMessages, deletedSessions)
}
