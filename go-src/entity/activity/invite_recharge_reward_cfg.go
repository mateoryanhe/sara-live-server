package activity

import (
	"xr-game-server/constants/db"
	"xr-game-server/core/migrate"
)

const (
	TbInviteRechargeRewardCfg db.TbName = "invite_recharge_reward_cfgs"
)

// InviteRechargeRewardCfg 邀请充值返还配置(CMS 管理,通常仅一条)
// 被邀请人完成账号首充后,在 ValidDays 内每次充值,邀请人获得到账金币 RewardPercent% 的返还。
type InviteRechargeRewardCfg struct {
	migrate.OneModel
	Enabled       bool    `gorm:"default:0;comment:活动开关" json:"enabled"`
	RewardPercent float64 `gorm:"type:decimal(6,2);default:5;comment:邀请人返还比例(%)" json:"rewardPercent"`
	ValidDays     int     `gorm:"default:30;comment:首充完成后有效天数" json:"validDays"`
}

func initInviteRechargeRewardCfg() {
	migrate.AutoMigrate(&InviteRechargeRewardCfg{})
}
