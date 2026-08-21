package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

type OnShelfShortVideoPriceTierReq struct {
	g.Meta `path:"/onShelfShortVideoPriceTier" method:"post" summary:"上架短视频价格挡位" tags:"短视频"`
	ID     uint64 `json:"id" v:"required#挡位ID不能为空" dc:"挡位ID"`
}

type OnShelfShortVideoPriceTierRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}

type OffShelfShortVideoPriceTierReq struct {
	g.Meta `path:"/offShelfShortVideoPriceTier" method:"post" summary:"下架短视频价格挡位" tags:"短视频"`
	ID     uint64 `json:"id" v:"required#挡位ID不能为空" dc:"挡位ID"`
}

type OffShelfShortVideoPriceTierRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}
