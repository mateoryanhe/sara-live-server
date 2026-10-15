package gameplatformdto

import "github.com/gogf/gf/v2/frame/g"

// AppGameStartReq App 获取游戏启动链接
type AppGameStartReq struct {
	g.Meta    `path:"/appGameStart" method:"post" summary:"App获取游戏启动链接" tags:"游戏"`
	GameCode  string `json:"gameCode" v:"required#游戏编码不能为空" dc:"游戏编码(对应第三方 gameId)"`
	RoomId    uint64 `json:"roomId" dc:"直播间ID:0=大厅玩游戏(清空关联);>0=该房间下注流水写入live_room_id"`
	UrlParams string `json:"urlParams" dc:"透传给厂商的 urlParams（原样上报）"`
}

// AppGameStartRes App 游戏启动链接响应
type AppGameStartRes struct {
	Link string `json:"link" dc:"游戏启动链接(有效期约2分钟,仅可使用一次)"`
}
