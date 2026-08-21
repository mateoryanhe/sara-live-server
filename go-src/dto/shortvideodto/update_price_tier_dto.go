package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

type UpdateShortVideoPriceTierReq struct {
	g.Meta `path:"/updateShortVideoPriceTier" method:"post" summary:"修改短视频价格挡位" tags:"短视频"`
	ID     uint64  `json:"id" v:"required#挡位ID不能为空" dc:"挡位ID"`
	Price  float64 `json:"price" dc:"钻石价格"`
}

type UpdateShortVideoPriceTierRes struct {
	Success bool `json:"success"`
}
