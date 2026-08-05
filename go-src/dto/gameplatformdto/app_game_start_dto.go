package gameplatformdto

import "github.com/gogf/gf/v2/frame/g"

// AppGameStartReq App 获取游戏启动链接
type AppGameStartReq struct {
	g.Meta   `path:"/appGameStart" method:"post" summary:"App获取游戏启动链接" tags:"游戏"`
	GameCode string `json:"gameCode" v:"required#游戏编码不能为空" dc:"游戏编码(对应第三方 gameId)"`
	Platform string `json:"platform" v:"required#平台编码不能为空" dc:"平台编码,如 pg / pp / wg"`
	Lang     string `json:"lang" dc:"语言,如 en / zh-CN,默认 en"`
}

// AppGameStartRes App 游戏启动链接响应
type AppGameStartRes struct {
	Link string `json:"link" dc:"游戏启动链接(有效期约2分钟,仅可使用一次)"`
}
