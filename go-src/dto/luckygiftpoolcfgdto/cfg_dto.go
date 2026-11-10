package luckygiftpoolcfgdto

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/dto/giftdto"
)

type GetLuckyGiftPoolCfgReq struct {
	g.Meta `path:"/getLuckyGiftPoolCfg" method:"post" summary:"查询幸运礼物奖池配置" tags:"幸运礼物奖池"`
}

type LuckyGiftPoolCfgItem struct {
	ID                string                              `json:"id"`
	PoolBalance       float64                             `json:"poolBalance"`
	ContributePercent float64                             `json:"contributePercent"`
	WinProbability    float64                             `json:"winProbability"`
	Tiers             []*giftdto.LuckyGiftMultiplierTierItem `json:"tiers"`
	CreatedAt         string                              `json:"createdAt"`
	UpdatedAt         string                              `json:"updatedAt"`
}

type GetLuckyGiftPoolCfgRes struct {
	Cfg *LuckyGiftPoolCfgItem `json:"cfg"`
}

type SaveLuckyGiftPoolCfgReq struct {
	g.Meta            `path:"/saveLuckyGiftPoolCfg" method:"post" summary:"保存幸运礼物奖池配置" tags:"幸运礼物奖池"`
	ID                string                              `json:"id" dc:"配置ID(字符串,避免前端 Number 丢 snowflake 精度);首次保存可空"`
	PoolBalance       float64                             `json:"poolBalance" v:"min:0#奖金池不能小于0" dc:"奖金池余额(默认10000000)"`
	ContributePercent float64                             `json:"contributePercent" v:"min:0|max:100#进池比例不能小于0|进池比例不能大于100" dc:"礼物金额进入奖池比例(%)"`
	WinProbability    float64                             `json:"winProbability" v:"min:0|max:100#中奖概率不能小于0|中奖概率不能大于100" dc:"中奖概率(%)"`
	Tiers             []*giftdto.LuckyGiftMultiplierTierItem `json:"tiers" v:"required#倍率挡位不能为空" dc:"5个倍率挡位"`
}

type SaveLuckyGiftPoolCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type ListLuckyGiftPoolPendingBalanceReq struct {
	g.Meta    `path:"/listPendingBalance" method:"post" summary:"幸运礼物奖池临时入账列表" tags:"幸运礼物奖池"`
	Status    *uint8 `json:"status" dc:"处理状态筛选(0未处理,1已处理,2处理中),不传=全部"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
}

type LuckyGiftPoolPendingBalanceItem struct {
	ID        string  `json:"id"`
	Amount    float64 `json:"amount"`
	Status    uint8   `json:"status"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}
