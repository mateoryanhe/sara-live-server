package luckygiftsettlementcfgdto

import "github.com/gogf/gf/v2/frame/g"

type GetLuckyGiftSettlementCfgReq struct {
	g.Meta `path:"/getLuckyGiftSettlementCfg" method:"post" summary:"查询幸运礼物结算配置" tags:"结算配置"`
}

type LuckyGiftSettlementCfgItem struct {
	ID                string  `json:"id"`
	SettlementPercent float64 `json:"settlementPercent"`
	CreatedAt         string  `json:"createdAt"`
	UpdatedAt         string  `json:"updatedAt"`
}

type GetLuckyGiftSettlementCfgRes struct {
	Cfg *LuckyGiftSettlementCfgItem `json:"cfg"`
}

type SaveLuckyGiftSettlementCfgReq struct {
	g.Meta            `path:"/saveLuckyGiftSettlementCfg" method:"post" summary:"保存幸运礼物结算配置" tags:"结算配置"`
	ID                uint64  `json:"id" dc:"配置ID,首次保存可为0"`
	SettlementPercent float64 `json:"settlementPercent" v:"min:0|max:100#结算比例不能小于0|结算比例不能大于100" dc:"幸运礼物结算比例(百分比,默认10)"`
}

type SaveLuckyGiftSettlementCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
