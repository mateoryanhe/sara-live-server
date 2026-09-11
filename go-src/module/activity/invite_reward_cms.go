package activity

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/activitydto"
	activityentity "xr-game-server/entity/activity"
	"xr-game-server/errercode"
)

func GetInviteRechargeRewardCfg(_ context.Context, _ *activitydto.GetInviteRechargeRewardCfgReq) (*activitydto.GetInviteRechargeRewardCfgRes, error) {
	return &activitydto.GetInviteRechargeRewardCfgRes{Cfg: toInviteRewardCfgItem(getInviteRewardCfgCache())}, nil
}

func SaveInviteRechargeRewardCfg(_ context.Context, req *activitydto.SaveInviteRechargeRewardCfgReq) (*activitydto.SaveInviteRechargeRewardCfgRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	percent := normalizeInviteRewardPercent(req.RewardPercent)
	days := normalizeInviteRewardDays(req.ValidDays)
	if percent < 0 || percent > 100 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if days < 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	existing := cfgdao.LoadInviteRechargeRewardCfg()
	row := &activityentity.InviteRechargeRewardCfg{
		Enabled:       req.Enabled,
		RewardPercent: percent,
		ValidDays:     days,
	}
	if req.ID > 0 {
		if existing == nil || existing.ID != req.ID {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		row.ID = req.ID
		row.CreatedAt = existing.CreatedAt
	} else if existing != nil {
		row.ID = existing.ID
		row.CreatedAt = existing.CreatedAt
	}
	row.UpdatedAt = time.Now()
	if row.CreatedAt.IsZero() {
		row.CreatedAt = row.UpdatedAt
	}
	if err := cfgdao.SaveInviteRechargeRewardCfg(row); err != nil {
		return nil, err
	}
	reloadInviteRewardCfgMemory()
	return &activitydto.SaveInviteRechargeRewardCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}
