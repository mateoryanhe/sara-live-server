package bannerdto

import "github.com/gogf/gf/v2/frame/g"

type OnShelfBannerReq struct {
	g.Meta `path:"/onShelfBanner" method:"post" summary:"上架首页Banner" tags:"首页Banner"`
	ID     uint64 `json:"id" v:"required#Banner ID不能为空" dc:"Banner ID"`
}

type OnShelfBannerRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}

type OffShelfBannerReq struct {
	g.Meta `path:"/offShelfBanner" method:"post" summary:"下架首页Banner" tags:"首页Banner"`
	ID     uint64 `json:"id" v:"required#Banner ID不能为空" dc:"Banner ID"`
}

type OffShelfBannerRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}
