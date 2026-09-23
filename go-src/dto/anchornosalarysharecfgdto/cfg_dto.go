package anchornosalarysharecfgdto

import "github.com/gogf/gf/v2/frame/g"

type GetAnchorNoSalaryShareCfgReq struct {
	g.Meta `path:"/getAnchorNoSalaryShareCfg" method:"post" summary:"查询无底薪主播提成配置" tags:"无底薪主播提成配置"`
}

type AnchorNoSalaryShareCfgItem struct {
	ID                       string  `json:"id"`
	AnchorSocialSharePercent float64 `json:"anchorSocialSharePercent"`
	GuildSocialSharePercent  float64 `json:"guildSocialSharePercent"`
	CreatedAt                string  `json:"createdAt"`
	UpdatedAt                string  `json:"updatedAt"`
}

type GetAnchorNoSalaryShareCfgRes struct {
	Cfg *AnchorNoSalaryShareCfgItem `json:"cfg"`
}

type SaveAnchorNoSalaryShareCfgReq struct {
	g.Meta                   `path:"/saveAnchorNoSalaryShareCfg" method:"post" summary:"保存无底薪主播提成配置" tags:"无底薪主播提成配置"`
	ID                       uint64  `json:"id" dc:"配置ID,新增传0"`
	AnchorSocialSharePercent float64 `json:"anchorSocialSharePercent" dc:"主播社交提成比(%)"`
	GuildSocialSharePercent  float64 `json:"guildSocialSharePercent" dc:"工会社交提成比(%)"`
}

type SaveAnchorNoSalaryShareCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
