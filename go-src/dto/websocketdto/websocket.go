package websocketdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type WebSocketAuthReq struct {
	g.Meta `path:"/auth" method:"post" description:"wss//:www.baidu.com:/ws?authId=10&token=kskksksks22 从登录成功接口获取" summary:"鉴权" tags:"推送" `
}

type WebSocketAuthRes struct {
	Cmd  int                 `json:"cmd" dc:"命令:1"`
	Data httpserver.AuthResp `json:"data" dc:"失败原因错误码"`
}
