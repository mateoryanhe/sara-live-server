package activitydto

import "github.com/gogf/gf/v2/frame/g"

type GetInviteRechargeRewardCfgReq struct {
	g.Meta `path:"/getInviteRechargeRewardCfg" method:"post" summary:"查询邀请充值返还配置" tags:"邀请充值返还"`
}

type InviteRechargeRewardCfgItem struct {
	ID            string  `json:"id"`
	Enabled       bool    `json:"enabled"`
	RewardPercent float64 `json:"rewardPercent"`
	ValidDays     int     `json:"validDays"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
}

type GetInviteRechargeRewardCfgRes struct {
	Cfg *InviteRechargeRewardCfgItem `json:"cfg"`
}

type SaveInviteRechargeRewardCfgReq struct {
	g.Meta        `path:"/saveInviteRechargeRewardCfg" method:"post" summary:"保存邀请充值返还配置" tags:"邀请充值返还"`
	ID            uint64  `json:"id" dc:"配置ID,首次保存可为0"`
	Enabled       bool    `json:"enabled" dc:"活动开关"`
	RewardPercent float64 `json:"rewardPercent" dc:"邀请人返还比例(%)"`
	ValidDays     int     `json:"validDays" dc:"首充完成后有效天数"`
}

type SaveInviteRechargeRewardCfgRes struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
