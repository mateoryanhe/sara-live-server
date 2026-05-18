package giftdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
)

// GiftListReq 分页查询礼物列表(CMS)
type GiftListReq struct {
	g.Meta `path:"/giftList" method:"post" summary:"获取礼物列表" tags:"礼物管理"`
	httpserver.CMSQueryReq
	Name         string `json:"name"     dc:"礼物名称(模糊匹配)"`
	Category     string `json:"category" dc:"分类"`
	StatusFilter int    `json:"statusFilter" dc:"状态过滤(0=全部, 1=只看下架, 2=只看上架)"`
}

type GiftListRes struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Animation   string `json:"animation"`
	Price       uint64 `json:"price"`
	Category    string `json:"category"`
	Sort        int    `json:"sort"`
	Status      uint8  `json:"status"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
