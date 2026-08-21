package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

type CreateShortVideoPriceTierReq struct {
	g.Meta `path:"/createShortVideoPriceTier" method:"post" summary:"创建短视频价格挡位" tags:"短视频"`
	Price  float64 `json:"price" dc:"钻石价格"`
}

type CreateShortVideoPriceTierRes struct {
	ID string `json:"id" dc:"挡位ID"`
}
