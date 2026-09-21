package userloginlocationdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/user"
)

var loginLocationCacheMgr *cache.RowCache[*entity.UserLoginLocation]

func Init() {
	loginLocationCacheMgr = cache.NewRowCache[*entity.UserLoginLocation]()
}

// GetByUserId 按 user_id 查询登录地域信息，先读进程缓存，未命中再查数据库。
// 数据库不存在时创建内存对象，字段修改通过 syndb 缓冲入库。
func GetByUserId(userId uint64) *entity.UserLoginLocation {
	if userId == 0 || loginLocationCacheMgr == nil {
		return nil
	}
	return loginLocationCacheMgr.MustGetRow(gctx.New(), userId, func(ctx context.Context) (*entity.UserLoginLocation, error) {
		var row *entity.UserLoginLocation
		err := g.Model(string(entity.TbUserLoginLocation)).Ctx(ctx).
			Where(string(entity.UserLoginLocationUserId), userId).
			Scan(&row)
		if err != nil {
			return nil, err
		}
		if row != nil && row.UserId > 0 {
			return row, nil
		}
		return entity.NewUserLoginLocation(userId), nil
	})
}

// GetByUserIds 批量查询登录地域信息，每个 user_id 都优先读取进程缓存，未命中项合并查询数据库。
func GetByUserIds(userIds []uint64) map[uint64]*entity.UserLoginLocation {
	result := make(map[uint64]*entity.UserLoginLocation, len(userIds))
	if len(userIds) == 0 || loginLocationCacheMgr == nil {
		return result
	}

	ctx := gctx.New()
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
		if row, ok := loginLocationCacheMgr.GetRowCached(ctx, userId); ok && row != nil {
			result[userId] = row
			continue
		}
		missing = append(missing, userId)
	}
	if len(missing) == 0 {
		return result
	}

	rows := make([]*entity.UserLoginLocation, 0, len(missing))
	if err := g.Model(string(entity.TbUserLoginLocation)).Ctx(ctx).
		WhereIn(string(entity.UserLoginLocationUserId), missing).
		Scan(&rows); err != nil {
		g.Log().Errorf(ctx, "load user login locations failed: %v", err)
		for _, userId := range missing {
			if row := GetByUserId(userId); row != nil {
				result[userId] = row
			}
		}
		return result
	}

	for _, row := range rows {
		if row == nil || row.UserId == 0 {
			continue
		}
		if cached, ok := loginLocationCacheMgr.GetRowCached(ctx, row.UserId); ok && cached != nil {
			result[row.UserId] = cached
			continue
		}
		loginLocationCacheMgr.PublishRow(ctx, row.UserId, row)
		result[row.UserId] = row
	}
	for _, userId := range missing {
		if result[userId] != nil {
			continue
		}
		if cached, ok := loginLocationCacheMgr.GetRowCached(ctx, userId); ok && cached != nil {
			result[userId] = cached
			continue
		}
		row := entity.NewUserLoginLocation(userId)
		loginLocationCacheMgr.PublishRow(ctx, userId, row)
		result[userId] = row
	}
	return result
}

// Publish 原地修改登录地域信息后刷新进程缓存。
func Publish(row *entity.UserLoginLocation) {
	if row == nil || row.UserId == 0 || loginLocationCacheMgr == nil {
		return
	}
	loginLocationCacheMgr.PublishRow(gctx.New(), row.UserId, row)
}

// GetFromMemory 仅查询进程缓存，未命中返回 nil。
func GetFromMemory(userId uint64) *entity.UserLoginLocation {
	if userId == 0 || loginLocationCacheMgr == nil {
		return nil
	}
	row, _ := loginLocationCacheMgr.GetRowCached(gctx.New(), userId)
	return row
}
