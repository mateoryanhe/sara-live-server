package effectivelivecfgdto

import "github.com/gogf/gf/v2/frame/g"

type GetEffectiveLiveCfgReq struct {
	g.Meta `path:"/getEffectiveLiveCfg" method:"post" summary:"查询有效直播时长配置" tags:"结算配置"`
}

type EffectiveLiveCfgItem struct {
	ID                       string `json:"id"`
	MinSessionMinutes        int    `json:"minSessionMinutes"`
	DailyAccumulatedMinutes  int    `json:"dailyAccumulatedMinutes"`
	CreatedAt                string `json:"createdAt"`
	UpdatedAt                string `json:"updatedAt"`
}

type GetEffectiveLiveCfgRes struct {
	Cfg *EffectiveLiveCfgItem `json:"cfg"`
}

type SaveEffectiveLiveCfgReq struct {
	g.Meta                  `path:"/saveEffectiveLiveCfg" method:"post" summary:"保存有效直播时长配置" tags:"结算配置"`
	ID                      uint64 `json:"id" dc:"配置ID,首次保存可为0"`
	MinSessionMinutes       int    `json:"minSessionMinutes" v:"required|min:1|max:1440#有效直播门槛不能为空|有效直播门槛至少为1分钟|有效直播门槛不能超过1440分钟" dc:"单场直播计入有效时长的门槛(分钟)"`
	DailyAccumulatedMinutes int    `json:"dailyAccumulatedMinutes" v:"min:0|max:1440#每天累计开播时长不能为负|每天累计开播时长不能超过1440分钟" dc:"每天累计开播时长门槛(分钟,0表示未配置)"`
}

type SaveEffectiveLiveCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
