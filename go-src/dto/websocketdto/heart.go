package websocketdto

import "github.com/gogf/gf/v2/frame/g"

type HeartReq struct {
	g.Meta `path:"/heart" method:"post" summary:"心跳" tags:"推送"`
}

type HeartRes struct {
	Cmd  int    `json:"cmd" dc:"命令:2"`
	Data uint64 `json:"data" dc:"时间戳"`
}
