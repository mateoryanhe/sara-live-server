package giftdto

import "github.com/gogf/gf/v2/frame/g"

// UpdateGiftReq 修改礼物
type UpdateGiftReq struct {
	g.Meta      `path:"/updateGift" method:"post" summary:"修改礼物" tags:"礼物管理"`
	ID          uint64 `json:"id"          v:"required#礼物ID不能为空" dc:"礼物ID"`
	Name        string `json:"name"        v:"required|length:1,64#礼物名称不能为空|名称长度需在1到64之间" dc:"礼物名称"`
	Icon        string `json:"icon"        v:"max-length:255#图标URL最长255字符" dc:"图标URL"`
	Animation   string `json:"animation"   v:"max-length:255#动画URL最长255字符" dc:"动画资源URL"`
	Price       uint64 `json:"price"       v:"min:0#价格不能为负" dc:"钻石单价"`
	Category    string `json:"category"    v:"max-length:32#分类最长32字符" dc:"分类"`
	Sort        int    `json:"sort"        dc:"排序值(越大越靠前)"`
	Description string `json:"description" v:"max-length:255#描述最长255字符" dc:"描述"`
}

type UpdateGiftRes struct {
	Success bool `json:"success"`
}
