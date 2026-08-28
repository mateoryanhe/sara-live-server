package userinfodao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/user"
)

var userRechargeCfgFirstRechargeCacheMgr *cache.ListCache[*entity.UserRechargeCfgFirstRecharge]

func initUserRechargeCfgFirstRechargeDao() {
	userRechargeCfgFirstRechargeCacheMgr = cache.NewListCache[*entity.UserRechargeCfgFirstRecharge]()
}

func getUserRechargeCfgFirstRechargeList(userId uint64) []*entity.UserRechargeCfgFirstRecharge {
	if userId == 0 || userRechargeCfgFirstRechargeCacheMgr == nil {
		return nil
	}
	return userRechargeCfgFirstRechargeCacheMgr.MustGetList(gctx.New(), userId, func(ctx context.Context) ([]*entity.UserRechargeCfgFirstRecharge, error) {
		rows := make([]*entity.UserRechargeCfgFirstRecharge, 0)
		_ = g.Model(string(entity.TbUserRechargeCfgFirstRecharge)).Ctx(ctx).Unscoped().
			Where(string(entity.UserRechargeCfgFirstRechargeUserId), userId).
			Order(string(entity.UserRechargeCfgFirstRechargeCfgId) + " asc").
			Scan(&rows)
		return rows, nil
	})
}

func publishUserRechargeCfgFirstRechargeList(userId uint64, list []*entity.UserRechargeCfgFirstRecharge) {
	if userId == 0 || userRechargeCfgFirstRechargeCacheMgr == nil {
		return
	}
	userRechargeCfgFirstRechargeCacheMgr.PublishList(gctx.New(), userId, list)
}

// IsRechargeCfgFirstRecharge 指定充值档位是否尚未首充(无记录视为未首充)
func IsRechargeCfgFirstRecharge(userId, cfgId uint64) bool {
	if userId == 0 || cfgId == 0 {
		return false
	}
	for _, row := range getUserRechargeCfgFirstRechargeList(userId) {
		if row != nil && row.CfgId == cfgId {
			return row.FirstRecharge
		}
	}
	return true
}

// HasAnyRechargeCfgFirstRecharge 是否存在任一已上架档位尚未首充
func HasAnyRechargeCfgFirstRecharge(userId uint64) bool {
	if userId == 0 {
		return false
	}
	done := make(map[uint64]struct{})
	for _, row := range getUserRechargeCfgFirstRechargeList(userId) {
		if row != nil && row.CfgId > 0 && !row.FirstRecharge {
			done[row.CfgId] = struct{}{}
		}
	}
	for _, cfg := range cfgdao.GetOnShelfRechargeCfg() {
		if cfg == nil || cfg.ID == 0 {
			continue
		}
		if _, ok := done[cfg.ID]; !ok {
			return true
		}
	}
	return false
}

// MarkRechargeCfgFirstRechargeDone 标记指定档位首充完成,返回本次是否为该档位首次到账
func MarkRechargeCfgFirstRechargeDone(userId, cfgId uint64) bool {
	if userId == 0 || cfgId == 0 {
		return false
	}
	list := getUserRechargeCfgFirstRechargeList(userId)
	var target *entity.UserRechargeCfgFirstRecharge
	for _, row := range list {
		if row != nil && row.CfgId == cfgId {
			target = row
			break
		}
	}
	if target != nil && !target.FirstRecharge {
		return false
	}
	newList := make([]*entity.UserRechargeCfgFirstRecharge, 0, len(list)+1)
	for _, row := range list {
		if row == nil {
			continue
		}
		if row.CfgId == cfgId {
			continue
		}
		newList = append(newList, row)
	}
	if target == nil {
		target = entity.NewUserRechargeCfgFirstRecharge(userId, cfgId)
	}
	target.SetFirstRecharge(false)
	newList = append(newList, target)
	publishUserRechargeCfgFirstRechargeList(userId, newList)
	return true
}

// PreloadUserRechargeCfgFirstRechargeToCache 批量预热用户档位首充列表缓存
func PreloadUserRechargeCfgFirstRechargeToCache(userIds []uint64) {
	if len(userIds) == 0 || userRechargeCfgFirstRechargeCacheMgr == nil {
		return
	}
	ctx := gctx.New()
	rows := make([]*entity.UserRechargeCfgFirstRecharge, 0)
	err := g.Model(string(entity.TbUserRechargeCfgFirstRecharge)).Ctx(ctx).Unscoped().
		WhereIn(string(entity.UserRechargeCfgFirstRechargeUserId), userIds).
		Order(string(entity.UserRechargeCfgFirstRechargeUserId) + " asc, " + string(entity.UserRechargeCfgFirstRechargeCfgId) + " asc").
		Scan(&rows)
	if err != nil {
		g.Log().Errorf(ctx, "preload user recharge cfg first recharge failed: %v", err)
		return
	}
	grouped := make(map[uint64][]*entity.UserRechargeCfgFirstRecharge, len(userIds))
	for _, row := range rows {
		if row == nil || row.UserId == 0 {
			continue
		}
		grouped[row.UserId] = append(grouped[row.UserId], row)
	}
	for userId, list := range grouped {
		userRechargeCfgFirstRechargeCacheMgr.PublishList(gctx.New(), userId, list)
	}
}
