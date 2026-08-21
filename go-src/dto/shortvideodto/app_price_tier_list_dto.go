package shortvideodto

import "github.com/gogf/gf/v2/frame/g"

type AppShortVideoPriceTierListReq struct {
	g.Meta `path:"/appShortVideoPriceTierList" method:"post" summary:"App查询短视频价格挡位列表(已上架)" tags:"短视频"`
}

type AppShortVideoPriceTierItem struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
}

type AppShortVideoPriceTierListRes struct {
	List []*AppShortVideoPriceTierItem `json:"list"`
}
