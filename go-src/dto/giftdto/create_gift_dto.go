package giftdto

import "github.com/gogf/gf/v2/frame/g"

// CreateGiftReq 创建礼物
type CreateGiftReq struct {
	g.Meta      `path:"/createGift" method:"post" summary:"创建礼物" tags:"礼物管理"`
	Name        string  `json:"name"        v:"required|length:1,64#礼物名称不能为空|名称长度需在1到64之间" dc:"礼物名称"`
	NameEn      string  `json:"nameEn"      v:"max-length:64#英文名称最长64字符" dc:"礼物名称(英文)"`
	NameEs      string  `json:"nameEs"      v:"max-length:64#西班牙语名称最长64字符" dc:"礼物名称(西班牙语)"`
	NamePt      string  `json:"namePt"      v:"max-length:64#葡萄牙语名称最长64字符" dc:"礼物名称(葡萄牙语)"`
	NameHi      string  `json:"nameHi"      v:"max-length:64#印地语名称最长64字符" dc:"礼物名称(印地语)"`
	Icon        string  `json:"icon"        v:"max-length:255#图标URL最长255字符" dc:"图标URL"`
	Animation   string  `json:"animation"   v:"max-length:255#动画URL最长255字符" dc:"动画资源URL"`
	Price       float64 `json:"price"       v:"min:0#价格不能为负" dc:"钻石单价"`
	Category    string  `json:"category"    v:"max-length:32#分类最长32字符" dc:"分类"`
	Sort        int     `json:"sort"          dc:"排序值(越大越靠前)"`
	Description string  `json:"description"   v:"max-length:255#描述最长255字符" dc:"描述"`
}

type CreateGiftRes struct {
	ID string `json:"id" dc:"礼物ID"`
}
