package platformanchorsettlementcfgdto

import "github.com/gogf/gf/v2/frame/g"

type GetPlatformAnchorSettlementCfgReq struct {
	g.Meta `path:"/getPlatformAnchorSettlementCfg" method:"post" summary:"查询平台主播结算配置" tags:"结算配置"`
}

type PlatformAnchorSettlementCfgItem struct {
	ID                   string  `json:"id"`
	MinimumSettlementUsd float64 `json:"minimumSettlementUsd"`
	CreatedAt            string  `json:"createdAt"`
	UpdatedAt            string  `json:"updatedAt"`
}

type GetPlatformAnchorSettlementCfgRes struct {
	Cfg *PlatformAnchorSettlementCfgItem `json:"cfg"`
}

type SavePlatformAnchorSettlementCfgReq struct {
	g.Meta               `path:"/savePlatformAnchorSettlementCfg" method:"post" summary:"保存平台主播结算配置" tags:"结算配置"`
	ID                   uint64  `json:"id" dc:"配置ID,首次保存可为0"`
	MinimumSettlementUsd float64 `json:"minimumSettlementUsd" v:"required|min:0.01|max:1000000#最低结算金额不能为空|最低结算金额必须大于0|最低结算金额不能超过1000000" dc:"平台主播生成结算单必须严格超过的最低金额(USD)"`
}

type SavePlatformAnchorSettlementCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
