package gamecfgdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// GameCfgListReq CMS分页查询游戏配置
type GameCfgListReq struct {
	g.Meta `path:"/gameCfgList" method:"post" summary:"获取游戏配置列表" tags:"游戏配置"`
	httpserver.CMSQueryReq
	Name         string `json:"name" dc:"游戏名称(模糊匹配)"`
	Code         string `json:"code" dc:"游戏编码(模糊匹配)"`
	StatusFilter int    `json:"statusFilter" dc:"状态过滤(0=全部,1=只看下架,2=只看上架)"`
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
	Status       uint8  `json:"status"`
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
	Status    uint8  `json:"status" v:"in:0,1#状态无效" dc:"状态(0下架,1上架)"`
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
	Status    uint8  `json:"status" v:"in:0,1#状态无效" dc:"状态(0下架,1上架)"`
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

// ===== App =====

// AppGameCfgListReq App端分页查询游戏列表(仅已上架,走缓存)
type AppGameCfgListReq struct {
	g.Meta   `path:"/appGameCfgList" method:"post" summary:"App分页查询游戏列表(已上架)" tags:"游戏配置"`
	Page     int `json:"page" dc:"页码(从1开始,默认1)"`
	PageSize int `json:"pageSize" dc:"每页数量(默认20,最大100)"`
}

// AppGameCfgItem App端游戏列表元素
type AppGameCfgItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	LiveCover string `json:"liveCover" dc:"直播间游戏封面URL"`
	Link      string `json:"link"`
	Sort      int    `json:"sort"`
}

// AppGameCfgListRes App端游戏分页列表响应
type AppGameCfgListRes struct {
	Total    int               `json:"total" dc:"总条数"`
	Page     int               `json:"page" dc:"当前页码"`
	PageSize int               `json:"pageSize" dc:"每页数量"`
	List     []*AppGameCfgItem `json:"list" dc:"游戏列表"`
}
