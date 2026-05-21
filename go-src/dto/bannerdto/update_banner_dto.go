package bannerdto

import "github.com/gogf/gf/v2/frame/g"

type UpdateBannerReq struct {
	g.Meta `path:"/updateBanner" method:"post" summary:"修改首页Banner" tags:"首页Banner"`
	ID     uint64 `json:"id" v:"required#Banner ID不能为空" dc:"Banner ID"`
	Title  string `json:"title" v:"required|length:1,64#标题不能为空|标题长度需在1到64之间" dc:"标题"`
	Image  string `json:"image" v:"max-length:255#图片资源名最长255字符" dc:"图片资源名"`
	Link   string `json:"link" v:"max-length:512#跳转链接最长512字符" dc:"跳转链接"`
	Sort   int    `json:"sort" dc:"排序值(越大越靠前)"`
}

type UpdateBannerRes struct {
	Success bool `json:"success"`
}
