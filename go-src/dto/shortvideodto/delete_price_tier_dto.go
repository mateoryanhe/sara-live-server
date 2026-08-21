package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

type DeleteShortVideoPriceTierReq struct {
	g.Meta `path:"/deleteShortVideoPriceTier" method:"post" summary:"删除短视频价格挡位" tags:"短视频"`
	ID     uint64 `json:"id" v:"required#挡位ID不能为空" dc:"挡位ID"`
}

type DeleteShortVideoPriceTierRes struct {
	Success bool `json:"success"`
}
