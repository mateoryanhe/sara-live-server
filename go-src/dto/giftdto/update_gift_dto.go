package giftdto

import "github.com/gogf/gf/v2/frame/g"

// UpdateGiftReq 修改礼物
type UpdateGiftReq struct {
	g.Meta      `path:"/updateGift" method:"post" summary:"修改礼物" tags:"礼物管理"`
	ID          uint64  `json:"id"          v:"required#礼物ID不能为空" dc:"礼物ID"`
	Name        string  `json:"name"        v:"required|length:1,64#礼物名称不能为空|名称长度需在1到64之间" dc:"礼物名称"`
	NameEn      string  `json:"nameEn"      v:"max-length:64#英文名称最长64字符" dc:"礼物名称(英文)"`
	NameEs      string  `json:"nameEs"      v:"max-length:64#西班牙语名称最长64字符" dc:"礼物名称(西班牙语)"`
	NamePt      string  `json:"namePt"      v:"max-length:64#葡萄牙语名称最长64字符" dc:"礼物名称(葡萄牙语)"`
	NameHi      string  `json:"nameHi"      v:"max-length:64#印地语名称最长64字符" dc:"礼物名称(印地语)"`
	NameId      string  `json:"nameId"      v:"max-length:64#印尼语名称最长64字符" dc:"礼物名称(印尼语)"`
	Icon        string  `json:"icon"        v:"max-length:255#图标URL最长255字符" dc:"图标URL"`
	Animation   string  `json:"animation"   v:"max-length:255#动画URL最长255字符" dc:"动画资源URL"`
	Price       float64 `json:"price"       v:"min:0#价格不能为负" dc:"钻石单价"`
	CategoryType uint8  `json:"categoryType" v:"in:0,1#礼物分类不合法" dc:"礼物分类(0普通,1幸运礼物)"`
	Sort         int     `json:"sort"          dc:"排序值(越大越靠前)"`
	Description string  `json:"description"   v:"max-length:255#描述最长255字符" dc:"描述"`
	LuckyGiftRtp                    float64                        `json:"luckyGiftRtp" v:"min:0|max:100#幸运礼物RTP不能小于0|幸运礼物RTP不能大于100" dc:"幸运礼物RTP(百分比,仅幸运礼物)"`
	LuckyGiftMultiplierTiers        []*LuckyGiftMultiplierTierItem `json:"luckyGiftMultiplierTiers" dc:"幸运礼物返奖倍率挡位(仅幸运礼物,至少1挡)"`
}

type UpdateGiftRes struct {
	Success bool `json:"success"`
}
