package event

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gutil"
	"xr-game-server/constants/common"
)

const (
	Online      Type = "Online"
	Offline     Type = "Offline"
	ClientEnter Type = "ClientEnter"
	ClientLeave Type = "ClientLeave"
)

type Handler func(val any)
type Type string

type OnlineData struct {
	RoleId uint64
}

func NewOnlineData(roleId uint64) *OnlineData {
	return &OnlineData{
		RoleId: roleId,
	}
}

type OfflineData struct {
	RoleId uint64
}

func NewOfflineData(roleId uint64) *OfflineData {
	return &OfflineData{
		RoleId: roleId,
	}
}

var eventHandlerMap = make(map[Type][]Handler)

// Sub 订阅事件
func Sub(eventType Type, handler Handler) {
	lst, ok := eventHandlerMap[eventType]
	if !ok {
		lst = make([]Handler, common.Zero)
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
			g.Log().Errorf(ctx, "event %s handler error %v", eventType, exception)
		})
	}
}
