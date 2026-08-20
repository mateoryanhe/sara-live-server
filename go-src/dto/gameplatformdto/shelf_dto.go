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
	Name     string `json:"name" dc:"游戏名称(模糊匹配,中/英文)"`
	Platform string `json:"platform" dc:"平台编码(模糊匹配)"`
}

// GameShelfListItem CMS 上架游戏项
type GameShelfListItem struct {
	ID               string `json:"id"`
	GameCode         string `json:"gameCode"`
	Name             string `json:"name"`
	NameEn           string `json:"nameEn"`
	Cover            string `json:"cover"`
	LiveGameName     string `json:"liveGameName"`
	LiveGameCover    string `json:"liveGameCover"`
	LiveGameCoverUrl string `json:"liveGameCoverUrl"`
	Platform         string `json:"platform"`
}

// AddGameShelfReq CMS 添加上架游戏
type AddGameShelfReq struct {
	g.Meta   `path:"/addGameShelf" method:"post" summary:"添加上架游戏" tags:"游戏平台配置"`
	GameCode string `json:"gameCode" v:"required#游戏编码不能为空" dc:"游戏编码"`
	Platform string `json:"platform" v:"required#平台编码不能为空" dc:"平台编码(第三方 platform 字段)"`
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

// BatchAddGameShelfItem CMS 批量上架单项
type BatchAddGameShelfItem struct {
	GameCode string `json:"gameCode" v:"required#游戏编码不能为空" dc:"游戏编码"`
	Platform string `json:"platform" v:"required#平台编码不能为空" dc:"平台编码(第三方 platform 字段)"`
}

// BatchAddGameShelfReq CMS 批量添加上架游戏
type BatchAddGameShelfReq struct {
	g.Meta `path:"/batchAddGameShelf" method:"post" summary:"批量添加上架游戏" tags:"游戏平台配置"`
	Items  []*BatchAddGameShelfItem `json:"items" v:"required#上架列表不能为空" dc:"上架游戏列表"`
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

// UpdateGameShelfReq CMS 更新上架游戏直播展示字段
type UpdateGameShelfReq struct {
	g.Meta        `path:"/updateGameShelf" method:"post" summary:"更新上架游戏直播展示" tags:"游戏平台配置"`
	GameCode      string `json:"gameCode" v:"required#游戏编码不能为空" dc:"游戏编码"`
	LiveGameName  string `json:"liveGameName" dc:"直播游戏名称(空则 App 回退英文名称)"`
	LiveGameCover string `json:"liveGameCover" dc:"直播游戏封面(相对路径或完整 URL, 空则 App 回退默认封面)"`
}

// UpdateGameShelfRes CMS 更新上架游戏直播展示结果
type UpdateGameShelfRes struct {
	Success bool `json:"success"`
}

// GetMultiplayerConfigUrlReq CMS 获取第三方自研游戏配置页链接
type GetMultiplayerConfigUrlReq struct {
	g.Meta   `path:"/getMultiplayerConfigUrl" method:"post" summary:"获取第三方自研游戏配置页链接" tags:"游戏平台配置"`
	GameCode string `json:"gameCode" v:"required#游戏编码不能为空" dc:"游戏编码(gameId)"`
	Platform string `json:"platform" dc:"平台编码(默认取上架记录, 否则 zy)"`
}

// GetMultiplayerConfigUrlRes CMS 第三方自研游戏配置页链接
type GetMultiplayerConfigUrlRes struct {
	ConfigUrl  string `json:"configUrl"`
	ExpireInMs int64  `json:"expireInMs"`
}

// CMSGameStartLinkReq CMS 代用户获取游戏启动链接
type CMSGameStartLinkReq struct {
	g.Meta   `path:"/cmsGameStartLink" method:"post" summary:"CMS代用户获取游戏启动链接" tags:"游戏平台配置"`
	UserId   uint64 `json:"userId,string" v:"required#用户ID不能为空" dc:"用户ID(作为第三方 ops)"`
	GameCode string `json:"gameCode" v:"required#游戏编码不能为空" dc:"游戏编码(对应第三方 gameId)"`
	Platform string `json:"platform" dc:"平台编码(默认取上架记录)"`
}

// CMSGameStartLinkRes CMS 游戏启动链接
type CMSGameStartLinkRes struct {
	Link string `json:"link" dc:"游戏启动链接(有效期约2分钟,仅可使用一次)"`
}
