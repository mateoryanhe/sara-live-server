package bannerdto

import "github.com/gogf/gf/v2/frame/g"

type DeleteBannerReq struct {
	g.Meta `path:"/deleteBanner" method:"post" summary:"删除首页Banner" tags:"首页Banner"`
	ID     uint64 `json:"id" v:"required#Banner ID不能为空" dc:"Banner ID"`
}

type DeleteBannerRes struct {
	Success bool `json:"success"`
}
