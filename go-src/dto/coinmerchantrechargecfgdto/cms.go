package coinmerchantrechargecfgdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// CoinMerchantRechargeCfgListReq CMS分页查询
type CoinMerchantRechargeCfgListReq struct {
	g.Meta `path:"/coinMerchantRechargeCfgList" method:"post" summary:"币商充值档位列表" tags:"币商充值档位"`
	httpserver.CMSQueryReq
	Name         string `json:"name" dc:"名称模糊"`
	StatusFilter int    `json:"statusFilter" dc:"0全部,1只看下架,2只看上架"`
}

// CoinMerchantRechargeCfgListRes 列表项
type CoinMerchantRechargeCfgListRes struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Gold      uint64  `json:"gold"`
	Status    uint8   `json:"status"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

// CreateCoinMerchantRechargeCfgReq 新建
type CreateCoinMerchantRechargeCfgReq struct {
	g.Meta `path:"/createCoinMerchantRechargeCfg" method:"post" summary:"新建币商充值档位" tags:"币商充值档位"`
	Name   string  `json:"name" v:"required|length:1,64#名称不能为空|名称长度1-64" dc:"档位名称"`
	Price  float64 `json:"price" v:"required|min:0.0001#价格不能为空|价格需大于0" dc:"USD价格"`
	Gold   uint64  `json:"gold" v:"required|min:1#到账金币不能为空|到账金币需大于0" dc:"到账金币"`
}

type CreateCoinMerchantRechargeCfgRes struct {
	ID string `json:"id"`
}

// UpdateCoinMerchantRechargeCfgReq 修改(不改上下架)
type UpdateCoinMerchantRechargeCfgReq struct {
	g.Meta `path:"/updateCoinMerchantRechargeCfg" method:"post" summary:"修改币商充值档位" tags:"币商充值档位"`
	ID     uint64  `json:"id" v:"required#ID不能为空"`
	Name   string  `json:"name" v:"required|length:1,64#名称不能为空|名称长度1-64"`
	Price  float64 `json:"price" v:"required|min:0.0001#价格不能为空|价格需大于0"`
	Gold   uint64  `json:"gold" v:"required|min:1#到账金币不能为空|到账金币需大于0"`
}

type UpdateCoinMerchantRechargeCfgRes struct {
	Success bool `json:"success"`
}

type DeleteCoinMerchantRechargeCfgReq struct {
	g.Meta `path:"/deleteCoinMerchantRechargeCfg" method:"post" summary:"删除币商充值档位" tags:"币商充值档位"`
	ID     uint64 `json:"id" v:"required#ID不能为空"`
}

type DeleteCoinMerchantRechargeCfgRes struct {
	Success bool `json:"success"`
}

type OnShelfCoinMerchantRechargeCfgReq struct {
	g.Meta `path:"/onShelfCoinMerchantRechargeCfg" method:"post" summary:"上架币商充值档位" tags:"币商充值档位"`
	ID     uint64 `json:"id" v:"required#ID不能为空"`
}

type OnShelfCoinMerchantRechargeCfgRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}

type OffShelfCoinMerchantRechargeCfgReq struct {
	g.Meta `path:"/offShelfCoinMerchantRechargeCfg" method:"post" summary:"下架币商充值档位" tags:"币商充值档位"`
	ID     uint64 `json:"id" v:"required#ID不能为空"`
}

type OffShelfCoinMerchantRechargeCfgRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}

// AppCoinMerchantRechargeCfgListReq App端查询币商充值档位(仅已上架)
type AppCoinMerchantRechargeCfgListReq struct {
	g.Meta `path:"/coinMerchantRechargeCfgListForApp" method:"post" summary:"App查询币商充值档位列表(已上架)" tags:"币商充值档位"`
}

// AppCoinMerchantRechargeCfgItem App/缓存条目
type AppCoinMerchantRechargeCfgItem struct {
	ID    uint64  `json:"id,string"`
	Name  string  `json:"name"`
	Price float64 `json:"price" dc:"现实货币价格(单位:USD)"`
	Gold  uint64  `json:"gold" dc:"到账金币数"`
}

// AppCoinMerchantRechargeCfgListRes App列表响应
type AppCoinMerchantRechargeCfgListRes struct {
	List []*AppCoinMerchantRechargeCfgItem `json:"list"`
}
