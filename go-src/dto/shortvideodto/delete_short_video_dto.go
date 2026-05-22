package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

type DeleteShortVideoReq struct {
	g.Meta `path:"/deleteShortVideo" method:"post" summary:"删除短视频" tags:"短视频"`
	ID     uint64 `json:"id" v:"required#短视频 ID不能为空" dc:"短视频 ID"`
}

type DeleteShortVideoRes struct {
	Success bool `json:"success"`
}
