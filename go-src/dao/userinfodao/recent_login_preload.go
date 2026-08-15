package userinfodao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	userentity "xr-game-server/entity/user"
)

// ListRecentLoginUserIds 按 last_login_time 倒序查询最近登录用户 ID
func ListRecentLoginUserIds(limit int) []uint64 {
	if limit <= 0 {
		return nil
	}
	rows := make([]*userentity.UserInfo, 0, limit)
	_ = g.Model(string(userentity.TbUserInfo)).Unscoped().
		Fields(string(db.IdName)).
		Where(string(userentity.UserInfoLastLoginTime) + " IS NOT NULL").
		Order(string(userentity.UserInfoLastLoginTime) + " desc").
		Limit(limit).
		Scan(&rows)
	userIds := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		userIds = append(userIds, row.ID)
	}
	return userIds
}

// PreloadRecentLoginUserInfos 批量加载并预热最近登录用户的 user_infos 缓存
func PreloadRecentLoginUserInfos(limit int) []uint64 {
	userIds := ListRecentLoginUserIds(limit)
	if len(userIds) == 0 {
		return nil
	}
	PreloadUserInfoToCache(loadUserInfosByUserIds(userIds))
	return userIds
}

func loadUserInfosByUserIds(userIds []uint64) []*userentity.UserInfo {
	if len(userIds) == 0 {
		return nil
	}
	rows := make([]*userentity.UserInfo, 0, len(userIds))
	ctx := gctx.New()
	err := g.Model(string(userentity.TbUserInfo)).Ctx(ctx).Unscoped().
		WhereIn(string(db.IdName), userIds).
		Scan(&rows)
	if err != nil {
		g.Log().Errorf(ctx, "preload user infos failed: %v", err)
		return nil
	}
	return rows
}

// PreloadUserInfoToCache 批量写入 user_infos 缓存
func PreloadUserInfoToCache(users []*userentity.UserInfo) {
	if len(users) == 0 || userInfoCacheMgr == nil {
		return
	}
	for _, user := range users {
		if user == nil || user.ID == 0 {
			continue
		}
		userInfoCacheMgr.FlushCache(user.ID, user)
		if user.ShareCode != "" && shareCodeUserIdCacheMgr != nil {
			shareCodeUserIdCacheMgr.FlushCache(user.ShareCode, user.ID)
		}
	}
}
