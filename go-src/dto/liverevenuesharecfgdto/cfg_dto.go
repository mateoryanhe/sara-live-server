package liverevenuesharecfgdto

import "github.com/gogf/gf/v2/frame/g"

type GetLiveRevenueShareCfgReq struct {
	g.Meta `path:"/getLiveRevenueShareCfg" method:"post" summary:"查询流水分佣配置" tags:"流水分佣配置"`
}

type LiveRevenueShareCfgItem struct {
	ID                 string  `json:"id"`
	AnchorSharePercent float64 `json:"anchorSharePercent"`
	GuildSharePercent  float64 `json:"guildSharePercent"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
}

type GetLiveRevenueShareCfgRes struct {
	Cfg *LiveRevenueShareCfgItem `json:"cfg"`
}

type SaveLiveRevenueShareCfgReq struct {
	g.Meta             `path:"/saveLiveRevenueShareCfg" method:"post" summary:"保存流水分佣配置" tags:"流水分佣配置"`
	ID                 uint64  `json:"id,string" dc:"配置ID,新增传0"`
	AnchorSharePercent float64 `json:"anchorSharePercent" v:"required#主播分佣比例不能为空"`
	GuildSharePercent  float64 `json:"guildSharePercent" v:"required#工会分佣比例不能为空"`
}

type SaveLiveRevenueShareCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
