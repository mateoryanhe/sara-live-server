package anchorvideopricetierdto

import "github.com/gogf/gf/v2/frame/g"

// AppPriceTierListReq App 查询主播相册短视频价格挡位(仅已上架,只读内存缓存)
type AppPriceTierListReq struct {
	g.Meta `path:"/priceTierList" method:"post" summary:"App查询主播相册短视频价格挡位(已上架,仅缓存)" tags:"主播相册视频"`
}

type AppPriceTierItem struct {
	ID    string  `json:"id" dc:"挡位ID"`
	Price float64 `json:"price" dc:"钻石价格"`
}

type AppPriceTierListRes struct {
	List []*AppPriceTierItem `json:"list" dc:"已上架价格挡位"`
}
