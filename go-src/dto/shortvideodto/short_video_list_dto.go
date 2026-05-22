package shortvideodto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type ShortVideoListReq struct {
	g.Meta `path:"/shortVideoList" method:"post" summary:"获取短视频列表" tags:"短视频"`
	httpserver.CMSQueryReq
	Title        string `json:"title" dc:"标题(模糊匹配)"`
	StatusFilter int    `json:"statusFilter" dc:"状态过滤(0=全部, 1=只看下架, 2=只看上架)"`
}

type ShortVideoListRes struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Video       string `json:"video" dc:"视频完整URL(列表展示)"`
	VideoName   string `json:"videoName" dc:"视频资源文件名(编辑保存用)"`
	Cover       string `json:"cover" dc:"封面完整URL(列表展示)"`
	CoverName   string `json:"coverName" dc:"封面资源文件名(编辑保存用)"`
	Sort        int    `json:"sort"`
	Status      uint8  `json:"status"`
	Description string `json:"description"`
	LikeCount   uint64 `json:"likeCount"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
