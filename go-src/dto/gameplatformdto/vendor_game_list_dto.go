package gameplatformdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// VendorGameListReq CMS 查询第三方游戏列表(内存缓存)
type VendorGameListReq struct {
	g.Meta `path:"/vendorGameList" method:"post" summary:"查询第三方游戏列表" tags:"游戏平台配置"`
	httpserver.CMSQueryReq
	GameCode string `json:"gameCode" dc:"游戏编码(模糊匹配)"`
	Name     string `json:"name" dc:"游戏名称(模糊匹配)"`
	Platform string `json:"platform" dc:"平台(模糊匹配)"`
	Category string `json:"category" dc:"分类(模糊匹配)"`
}

// VendorGameListItem CMS 游戏列表项
type VendorGameListItem struct {
	GameCode string `json:"gameCode"`
	Name     string `json:"name"`
	NameEn   string `json:"nameEn"`
	Category string `json:"category"`
	Cover    string `json:"cover"`
	Platform string `json:"platform"`
	OnShelf  bool   `json:"onShelf"`
}

// ReloadVendorGameCacheReq CMS 从第三方重新拉取游戏列表
type ReloadVendorGameCacheReq struct {
	g.Meta `path:"/reloadVendorGameCache" method:"post" summary:"从第三方重新拉取游戏列表" tags:"游戏平台配置"`
}

// ReloadVendorGameCacheRes CMS 重新拉取结果
type ReloadVendorGameCacheRes struct {
	Success bool `json:"success"`
	Count   int  `json:"count"`
}
