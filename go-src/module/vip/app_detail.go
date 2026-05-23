package vip

import (
	"context"
	"xr-game-server/core/httpserver"
	coremath "xr-game-server/core/ma
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/vipcfgdto"
	"xr-game-server/dto/vipdto"
)

// GetAppVipDetail App端查询当前登录用户VIP详情
func GetAppVipDetail(ctx context.Context, _ *vipdto.AppVipDetailReq) (*vipdto.AppVipDetailRes, error) {
	userId := httpserver.GetAuthId(ctx)
	user := userinfodao.GetUserInfoByUserId(userId)
	stat := userinfodao.GetUserCumulativeStatByUserId(userId)

	totalRecharge := coremath.RoundFloat64(stat.TotalRecharge)
	res := &vipdto.AppVipDetailRes{
		VipLevel:      user.VipLevel,
		TotalRecharge: totalRecharge,
	}

	if user.VipLevel > 0 {
		if cfg := GetVipCfgFromMemoryByLevel(user.VipLevel); cfg != nil {
			res.CurrentCfg = cfg
			res.LevelName = cfg.LevelName
		}
	}

	nextCfg := findNextEnabledVipCfg(user.VipLevel)
	if nextCfg == nil {
		res.IsMaxLevel = true
		return res, nil
	}

	res.NextLevel = nextCfg.Level
	res.NextCfg = nextCfg
	res.NextUpgradeLimit = nextCfg.UpgradeRechargeLimit
	if totalRecharge < nextCfg.UpgradeRechargeLimit {
		res.RechargeToNextLevel = coremath.SubFloat64(nextCfg.UpgradeRechargeLimit, totalRecharge)
	}
	return res, nil
}

func findNextEnabledVipCfg(currentLevel uint32) *vipcfgdto.AppVipCfgItem {
	var next *vipcfgdto.AppVipCfgItem
	for _, cfg := range GetAllVipCfgFromMemory() {
		if cfg.Level <= currentLevel {
			continue
		}
		if next == nil || cfg.Level < next.Level {
			next = cfg
		}
	}
	return next
}
