package giftdto

import "github.com/gogf/gf/v2/frame/g"

type LuckyGiftMultiplierTierItem struct {
	ID         string `json:"id" dc:"挡位ID"`
	Multiplier uint32 `json:"multiplier" dc:"返奖倍率"`
	Weight     uint32 `json:"weight" dc:"抽奖权重"`
	Type       string `json:"type" dc:"挡位类型(low/medium/high)"`
}

type GetGiftLuckyCfgReq struct {
	g.Meta `path:"/getGiftLuckyCfg" method:"post" summary:"查询单个幸运礼物RTP与倍率挡位" tags:"礼物管理"`
	GiftID uint64 `json:"giftId" v:"required#礼物ID不能为空" dc:"礼物ID"`
}

type GetGiftLuckyCfgRes struct {
	LuckyGiftRtp float64                       `json:"luckyGiftRtp"`
	Tiers        []*LuckyGiftMultiplierTierItem `json:"tiers"`
}
