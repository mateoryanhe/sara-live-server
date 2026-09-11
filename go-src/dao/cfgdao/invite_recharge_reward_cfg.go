package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/activity"
)

func LoadInviteRechargeRewardCfg() *activity.InviteRechargeRewardCfg {
	var row activity.InviteRechargeRewardCfg
	if err := g.DB().Model(string(activity.TbInviteRechargeRewardCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SaveInviteRechargeRewardCfg(row *activity.InviteRechargeRewardCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(activity.TbInviteRechargeRewardCfg)).Save(row)
	return err
}
