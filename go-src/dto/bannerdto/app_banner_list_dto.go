package bannerdto

import "github.com/gogf/gf/v2/frame/g"

type AppBannerListReq struct {
	g.Meta `path:"/appBannerList" method:"post" summary:"App查询首页Banner(已上架)" tags:"首页Banner"`
}

type AppBannerItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Image     string `json:"image"`
	Link      string `json:"link"`
	Direction uint8  `json:"direction"`
	Sort      int    `json:"sort"`
}

type AppBannerListRes struct {
	List []*AppBannerItem `json:"list"`
}
