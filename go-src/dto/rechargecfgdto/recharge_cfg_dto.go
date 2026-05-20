package rechargecfgdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// ===== CMS =====

// RechargeCfgListReq CMS分页查询充值配置
type RechargeCfgListReq struct {
	g.Meta `path:"/rechargeCfgList" method:"post" summary:"获取充值配置列表" tags:"充值配置"`
	httpserver.CMSQueryReq
	Name         string `json:"name"         dc:"名称(模糊匹配)"`
	StatusFilter int    `json:"statusFilter" dc:"状态过滤(0=全部, 1=只看下架, 2=只看上架)"`
}

// RechargeCfgListRes 单条列表项
type RechargeCfgListRes struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Icon         string `json:"icon"`
	Diamond      uint64 `json:"diamond"`
	ExtraDiamond uint64 `json:"extraDiamond"`
	Price        uint64 `json:"price"`
	Currency     string `json:"currency"`
	ProductId    string `json:"productId"`
	Sort         int    `json:"sort"`
	Status       uint8  `json:"status"`
	Description  string `json:"description"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// CreateRechargeCfgReq 创建充值配置
type CreateRechargeCfgReq struct {
	g.Meta       `path:"/createRechargeCfg" method:"post" summary:"创建充值配置" tags:"充值配置"`
	Name         string `json:"name"         v:"required|length:1,64#名称不能为空|名称长度需在1到64之间" dc:"档位名称"`
	Icon         string `json:"icon"         v:"max-length:255#图标URL最长255字符" dc:"图标URL"`
	Diamond      uint64 `json:"diamond"      v:"required|min:1#基础到账钻石数不能为空|基础到账钻石数需大于0" dc:"基础到账钻石数"`
	ExtraDiamond uint64 `json:"extraDiamond" dc:"额外赠送钻石数"`
	Price        uint64 `json:"price"        v:"required|min:1#价格不能为空|价格需大于0" dc:"现实货币价格(单位:分)"`
	Currency     string `json:"currency"     v:"max-length:8#货币代码最长8字符" dc:"货币(CNY/USD等),默认CNY"`
	ProductId    string `json:"productId"    v:"max-length:64#商品ID最长64字符" dc:"第三方商品SKU"`
	Sort         int    `json:"sort"         dc:"排序值(越大越靠前)"`
	Description  string `json:"description"  v:"max-length:255#描述最长255字符" dc:"描述"`
}

type CreateRechargeCfgRes struct {
	ID string `json:"id"`
}

// UpdateRechargeCfgReq 更新充值配置(不修改上下架状态)
type UpdateRechargeCfgReq struct {
	g.Meta       `path:"/updateRechargeCfg" method:"post" summary:"修改充值配置" tags:"充值配置"`
	ID           uint64 `json:"id"           v:"required#ID不能为空" dc:"档位ID"`
	Name         string `json:"name"         v:"required|length:1,64#名称不能为空|名称长度需在1到64之间" dc:"档位名称"`
	Icon         string `json:"icon"         v:"max-length:255#图标URL最长255字符" dc:"图标URL"`
	Diamond      uint64 `json:"diamond"      v:"required|min:1#基础到账钻石数不能为空|基础到账钻石数需大于0" dc:"基础到账钻石数"`
	ExtraDiamond uint64 `json:"extraDiamond" dc:"额外赠送钻石数"`
	Price        uint64 `json:"price"        v:"required|min:1#价格不能为空|价格需大于0" dc:"现实货币价格(单位:分)"`
	Currency     string `json:"currency"     v:"max-length:8#货币代码最长8字符" dc:"货币(CNY/USD等)"`
	ProductId    string `json:"productId"    v:"max-length:64#商品ID最长64字符" dc:"第三方商品SKU"`
	Sort         int    `json:"sort"         dc:"排序值"`
	Description  string `json:"description"  v:"max-length:255#描述最长255字符" dc:"描述"`
}

type UpdateRechargeCfgRes struct {
	Success bool `json:"success"`
}

// DeleteRechargeCfgReq 删除充值配置
type DeleteRechargeCfgReq struct {
	g.Meta `path:"/deleteRechargeCfg" method:"post" summary:"删除充值配置" tags:"充值配置"`
	ID     uint64 `json:"id" v:"required#ID不能为空" dc:"档位ID"`
}

type DeleteRechargeCfgRes struct {
	Success bool `json:"success"`
}

// OnShelfRechargeCfgReq 上架
type OnShelfRechargeCfgReq struct {
	g.Meta `path:"/onShelfRechargeCfg" method:"post" summary:"上架充值配置" tags:"充值配置"`
	ID     uint64 `json:"id" v:"required#ID不能为空" dc:"档位ID"`
}

type OnShelfRechargeCfgRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}

// OffShelfRechargeCfgReq 下架
type OffShelfRechargeCfgReq struct {
	g.Meta `path:"/offShelfRechargeCfg" method:"post" summary:"下架充值配置" tags:"充值配置"`
	ID     uint64 `json:"id" v:"required#ID不能为空" dc:"档位ID"`
}

type OffShelfRechargeCfgRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}

// ===== App =====

// AppRechargeCfgListReq App端查询充值配置(仅返回已上架)
type AppRechargeCfgListReq struct {
	g.Meta `path:"/rechargeCfgListForApp" method:"post" summary:"App查询充值配置列表(已上架)" tags:"充值配置"`
}

// AppRechargeCfgItem App端单条
type AppRechargeCfgItem struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Icon         string `json:"icon"`
	Diamond      uint64 `json:"diamond"`
	ExtraDiamond uint64 `json:"extraDiamond"`
	Price        uint64 `json:"price"`
	Currency     string `json:"currency"`
	ProductId    string `json:"productId"`
	Sort         int    `json:"sort"`
	Description  string `json:"description"`
}

type AppRechargeCfgListRes struct {
	List []*AppRechargeCfgItem `json:"list"`
}
