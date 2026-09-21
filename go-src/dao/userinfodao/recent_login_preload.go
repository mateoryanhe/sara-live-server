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
		userInfoCacheMgr.PublishRow(gctx.New(), user.ID, user)
	}
}

// EnsureUserInfosCached 批量确保用户基础信息已进入进程缓存。
// 已命中的用户直接返回；所有未命中的用户合并为一次数据库查询并回填缓存。
func EnsureUserInfosCached(userIds []uint64) map[uint64]*userentity.UserInfo {
	result := make(map[uint64]*userentity.UserInfo, len(userIds))
	if len(userIds) == 0 || userInfoCacheMgr == nil {
		return result
	}

	missing := make([]uint64, 0, len(userIds))
	seen := make(map[uint64]struct{}, len(userIds))
	for _, userId := range userIds {
		if userId == 0 {
			continue
		}
		if _, ok := seen[userId]; ok {
			continue
		}
		seen[userId] = struct{}{}
		if user := GetUserInfoFromMemory(userId); user != nil {
			result[userId] = user
			continue
		}
		missing = append(missing, userId)
	}

	loaded := loadUserInfosByUserIds(missing)
	PreloadUserInfoToCache(loaded)
	for _, user := range loaded {
		if user == nil || user.ID == 0 {
			continue
		}
		result[user.ID] = user
	}
	return result
}
