package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

type OnShelfShortVideoReq struct {
	g.Meta `path:"/onShelfShortVideo" method:"post" summary:"上架短视频" tags:"短视频"`
	ID     uint64 `json:"id" v:"required#短视频 ID不能为空" dc:"短视频 ID"`
}

type OnShelfShortVideoRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}

type OffShelfShortVideoReq struct {
	g.Meta `path:"/offShelfShortVideo" method:"post" summary:"下架短视频" tags:"短视频"`
	ID     uint64 `json:"id" v:"required#短视频 ID不能为空" dc:"短视频 ID"`
}

type OffShelfShortVideoRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}
