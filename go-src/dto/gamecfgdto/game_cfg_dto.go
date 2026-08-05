package gamecfgdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// GameCfgListReq CMS分页查询游戏配置
type GameCfgListReq struct {
	g.Meta `path:"/gameCfgList" method:"post" summary:"获取游戏配置列表" tags:"游戏配置"`
	httpserver.CMSQueryReq
	Name string `json:"name" dc:"游戏名称(模糊匹配)"`
	Code string `json:"code" dc:"游戏编码(模糊匹配)"`
}

// GameCfgListRes 列表项
type GameCfgListRes struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	LiveCover    string `json:"liveCover"`
	LiveCoverUrl string `json:"liveCoverUrl"`
	Link         string `json:"link"`
	Sort         int    `json:"sort"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// CreateGameCfgReq 创建游戏配置
type CreateGameCfgReq struct {
	g.Meta    `path:"/createGameCfg" method:"post" summary:"创建游戏配置" tags:"游戏配置"`
	Name      string `json:"name" v:"required|length:1,64#游戏名称不能为空|游戏名称长度需在1到64之间" dc:"游戏名称"`
	Code      string `json:"code" v:"required|length:1,64#游戏编码不能为空|游戏编码长度需在1到64之间" dc:"游戏编码"`
	LiveCover string `json:"liveCover" v:"required#请上传直播间游戏封面" dc:"直播间游戏封面资源名"`
	Link      string `json:"link" dc:"跳转链接"`
	Sort      int    `json:"sort" dc:"排序值(越大越靠前)"`
}

type CreateGameCfgRes struct {
	ID string `json:"id"`
}

// UpdateGameCfgReq 更新游戏配置
type UpdateGameCfgReq struct {
	g.Meta    `path:"/updateGameCfg" method:"post" summary:"修改游戏配置" tags:"游戏配置"`
	ID        uint64 `json:"id" v:"required#ID不能为空" dc:"配置ID"`
	Name      string `json:"name" v:"required|length:1,64#游戏名称不能为空|游戏名称长度需在1到64之间" dc:"游戏名称"`
	Code      string `json:"code" v:"required|length:1,64#游戏编码不能为空|游戏编码长度需在1到64之间" dc:"游戏编码"`
	LiveCover string `json:"liveCover" v:"required#请上传直播间游戏封面" dc:"直播间游戏封面资源名"`
	Link      string `json:"link" dc:"跳转链接"`
	Sort      int    `json:"sort" dc:"排序值(越大越靠前)"`
}

type UpdateGameCfgRes struct {
	Success bool `json:"success"`
}

// DeleteGameCfgReq 删除游戏配置
type DeleteGameCfgReq struct {
	g.Meta `path:"/deleteGameCfg" method:"post" summary:"删除游戏配置" tags:"游戏配置"`
	ID     uint64 `json:"id" v:"required#ID不能为空" dc:"配置ID"`
}

type DeleteGameCfgRes struct {
	Success bool `json:"success"`
}
