package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

type AppDeleteShortVideoReq struct {
	g.Meta `path:"/appDeleteShortVideo" method:"post" summary:"App删除未审核通过的短视频" tags:"短视频"`
	ID     uint64 `json:"id,string" v:"required#短视频ID不能为空" dc:"短视频ID"`
}

type AppDeleteShortVideoRes struct {
	Success bool `json:"success"`
}
