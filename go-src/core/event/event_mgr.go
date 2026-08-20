package event

import (
	"context"
	"xr-game-server/core/xrlog"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gutil"
)

const (
	Online                 Type = "Online"
	Offline                Type = "Offline"
	ClientEnter            Type = "ClientEnter"
	ClientLeave            Type = "ClientLeave"
	AppToken               Type = "AppToken"
	CmsToken               Type = "CmsToken"
	PrepareRestart         Type = "PrepareRestart"
	LiveRoomAudienceJoined Type = "LiveRoomAudienceJoined"
)

type Handler func(val any)
type Type string

var eventHandlerMap = make(map[Type][]Handler)

// Sub 订阅事件
func Sub(eventType Type, handler Handler) {
	lst, ok := eventHandlerMap[eventType]
	if !ok {
		lst = make([]Handler, 0)
	}
	eventHandlerMap[eventType] = append(lst, handler)
}

// Pub 发布事件
func Pub(eventType Type, eventVal any) {
	handlers, ok := eventHandlerMap[eventType]
	if !ok {
		return
	}
	if len(handlers) == 0 {
		return
	}
	for _, handler := range handlers {
		gutil.TryCatch(gctx.New(), func(ctx context.Context) {
			handler(eventVal)
		}, func(ctx context.Context, exception error) {
			xrlog.ErrorWithErr(ctx, "Event", "handler="+string(eventType), exception)
		})
	}
}
