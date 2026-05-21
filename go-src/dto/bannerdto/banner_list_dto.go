package bannerdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type BannerListReq struct {
	g.Meta `path:"/bannerList" method:"post" summary:"获取首页Banner列表" tags:"首页Banner"`
	httpserver.CMSQueryReq
	Title        string `json:"title" dc:"标题(模糊匹配)"`
	StatusFilter int    `json:"statusFilter" dc:"状态过滤(0=全部, 1=只看下架, 2=只看上架)"`
}

type BannerListRes struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Image     string `json:"image" dc:"图片完整URL(列表展示)"`
	ImageName string `json:"imageName" dc:"图片资源文件名(编辑保存用)"`
	Link      string `json:"link"`
	Sort      int    `json:"sort"`
	Status    uint8  `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
