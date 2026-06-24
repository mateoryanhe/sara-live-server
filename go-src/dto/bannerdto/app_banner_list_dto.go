package bannerdto

import "github.com/gogf/gf/v2/frame/g"

type AppBannerListReq struct {
	g.Meta `path:"/appBannerList" method:"post" summary:"App查询Banner(已上架)" tags:"首页Banner"`
	Scene  uint8 `json:"scene" v:"in:0,1,2#展示场景无效" dc:"展示场景(0/1首页,2直播间),默认首页"`
}

type AppBannerItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Image     string `json:"image"`
	Link      string `json:"link"`
	Scene     uint8  `json:"scene"`
	Direction uint8  `json:"direction"`
	Sort      int    `json:"sort"`
}

type AppBannerListRes struct {
	List []*AppBannerItem `json:"list"`
}
