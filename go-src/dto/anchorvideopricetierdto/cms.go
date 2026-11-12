package anchorvideopricetierdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

type PriceTierListReq struct {
	g.Meta `path:"/anchorVideoPriceTierList" method:"post" summary:"获取主播相册短视频价格挡位列表" tags:"主播相册短视频价格"`
	httpserver.CMSQueryReq
	StatusFilter int `json:"statusFilter" dc:"状态过滤(0=全部,1=只看下架,2=只看上架)"`
}

type PriceTierListItem struct {
	ID        string  `json:"id"`
	Price     float64 `json:"price"`
	Status    uint8   `json:"status"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

type CreatePriceTierReq struct {
	g.Meta `path:"/createAnchorVideoPriceTier" method:"post" summary:"创建主播相册短视频价格挡位" tags:"主播相册短视频价格"`
	Price  float64 `json:"price" dc:"钻石价格"`
}

type CreatePriceTierRes struct {
	ID string `json:"id" dc:"挡位ID"`
}

type UpdatePriceTierReq struct {
	g.Meta `path:"/updateAnchorVideoPriceTier" method:"post" summary:"修改主播相册短视频价格挡位" tags:"主播相册短视频价格"`
	ID     uint64  `json:"id" v:"required#挡位ID不能为空" dc:"挡位ID"`
	Price  float64 `json:"price" dc:"钻石价格"`
}

type UpdatePriceTierRes struct {
	Success bool `json:"success"`
}

type DeletePriceTierReq struct {
	g.Meta `path:"/deleteAnchorVideoPriceTier" method:"post" summary:"删除主播相册短视频价格挡位" tags:"主播相册短视频价格"`
	ID     uint64 `json:"id" v:"required#挡位ID不能为空" dc:"挡位ID"`
}

type DeletePriceTierRes struct {
	Success bool `json:"success"`
}

type OnShelfPriceTierReq struct {
	g.Meta `path:"/onShelfAnchorVideoPriceTier" method:"post" summary:"上架主播相册短视频价格挡位" tags:"主播相册短视频价格"`
	ID     uint64 `json:"id" v:"required#挡位ID不能为空" dc:"挡位ID"`
}

type OnShelfPriceTierRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}

type OffShelfPriceTierReq struct {
	g.Meta `path:"/offShelfAnchorVideoPriceTier" method:"post" summary:"下架主播相册短视频价格挡位" tags:"主播相册短视频价格"`
	ID     uint64 `json:"id" v:"required#挡位ID不能为空" dc:"挡位ID"`
}

type OffShelfPriceTierRes struct {
	Success bool  `json:"success"`
	Status  uint8 `json:"status"`
}
