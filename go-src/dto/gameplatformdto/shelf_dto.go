package gameplatformdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// GameShelfListReq CMS 查询上架游戏列表(内存)
type GameShelfListReq struct {
	g.Meta `path:"/gameShelfList" method:"post" summary:"查询上架游戏列表" tags:"游戏平台配置"`
	httpserver.CMSQueryReq
	GameCode string `json:"gameCode" dc:"游戏编码(模糊匹配)"`
}

// GameShelfListItem CMS 上架游戏项
type GameShelfListItem struct {
	ID       string `json:"id"`
	GameCode string `json:"gameCode"`
	Platform string `json:"platform"`
}

// AddGameShelfReq CMS 添加上架游戏
type AddGameShelfReq struct {
	g.Meta   `path:"/addGameShelf" method:"post" summary:"添加上架游戏" tags:"游戏平台配置"`
	GameCode string `json:"gameCode" v:"required#游戏编码不能为空" dc:"游戏编码"`
}

// AddGameShelfRes CMS 添加上架游戏结果
type AddGameShelfRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

// DeleteGameShelfReq CMS 删除上架游戏
type DeleteGameShelfReq struct {
	g.Meta   `path:"/deleteGameShelf" method:"post" summary:"删除上架游戏" tags:"游戏平台配置"`
	ID       uint64 `json:"id" dc:"记录ID"`
	GameCode string `json:"gameCode" dc:"游戏编码"`
}

// DeleteGameShelfRes CMS 删除上架游戏结果
type DeleteGameShelfRes struct {
	Success bool `json:"success"`
}

// BatchAddGameShelfReq CMS 批量添加上架游戏
type BatchAddGameShelfReq struct {
	g.Meta    `path:"/batchAddGameShelf" method:"post" summary:"批量添加上架游戏" tags:"游戏平台配置"`
	GameCodes []string `json:"gameCodes" v:"required#游戏编码列表不能为空" dc:"游戏编码列表"`
}

// BatchAddGameShelfRes CMS 批量添加上架游戏结果
type BatchAddGameShelfRes struct {
	Success      bool `json:"success"`
	SuccessCount int  `json:"successCount"`
	SkipCount    int  `json:"skipCount"`
}

// BatchDeleteGameShelfReq CMS 批量删除上架游戏
type BatchDeleteGameShelfReq struct {
	g.Meta    `path:"/batchDeleteGameShelf" method:"post" summary:"批量删除上架游戏" tags:"游戏平台配置"`
	GameCodes []string `json:"gameCodes" v:"required#游戏编码列表不能为空" dc:"游戏编码列表"`
}

// BatchDeleteGameShelfRes CMS 批量删除上架游戏结果
type BatchDeleteGameShelfRes struct {
	Success      bool `json:"success"`
	SuccessCount int  `json:"successCount"`
}
