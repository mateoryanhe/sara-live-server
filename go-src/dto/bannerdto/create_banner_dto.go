package bannerdto

import "github.com/gogf/gf/v2/frame/g"

type CreateBannerReq struct {
	g.Meta `path:"/createBanner" method:"post" summary:"创建首页Banner" tags:"首页Banner"`
	Title  string `json:"title" v:"required|length:1,64#标题不能为空|标题长度需在1到64之间" dc:"标题"`
	Image  string `json:"image" v:"max-length:255#图片资源名最长255字符" dc:"图片资源名"`
	Link   string `json:"link" v:"max-length:512#跳转链接最长512字符" dc:"跳转链接"`
	Sort   int    `json:"sort" dc:"排序值(越大越靠前)"`
}

type CreateBannerRes struct {
	ID string `json:"id" dc:"Banner ID"`
}
