package websocketdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type WebSocketReq struct {
	g.Meta `path:"/base" method:"post" summary:"推送协议" tags:"推送"`
}

type WebSocketRes struct {
	Data []*httpserver.PushResp `json:"data"`
}
