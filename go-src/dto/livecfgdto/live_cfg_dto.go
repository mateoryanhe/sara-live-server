package livecfgdto

import "github.com/gogf/gf/v2/frame/g"

type GetLiveCfgReq struct {
	g.Meta `path:"/getLiveCfg" method:"post" summary:"查询直播配置" tags:"直播配置"`
}

type LiveCfgItem struct {
	ID               string  `json:"id"`
	PaidDanmakuPrice float64 `json:"paidDanmakuPrice"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
}

type GetLiveCfgRes struct {
	Cfg *LiveCfgItem `json:"cfg"`
}

type SaveLiveCfgReq struct {
	g.Meta           `path:"/saveLiveCfg" method:"post" summary:"保存直播配置" tags:"直播配置"`
	ID               uint64  `json:"id" dc:"配置ID,首次保存可为0"`
	PaidDanmakuPrice float64 `json:"paidDanmakuPrice" v:"required|min:0#付费弹幕价格不能为空|付费弹幕价格不能小于0"`
}

type SaveLiveCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
