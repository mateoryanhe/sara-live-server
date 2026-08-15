package calldao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity/call"
)

var userCacheMgr *cache.CacheMgr

// GetUserById 按主键(userId)查询通话用户(走缓存)
func GetUserById(userId uint64) *entity.CallUser {
	if userId == 0 || userCacheMgr == nil {
		return nil
	}
	v := userCacheMgr.GetData(userId, func(ctx context.Context) (value interface{}, err error) {
		var ret entity.CallUser
		err = g.DB().Model(string(entity.TbCallUser)).WherePri(userId).Scan(&ret)
		if err != nil || ret.ID == 0 {
			return nil, err
		}
		return &ret, nil
	})
	if v == nil {
		return nil
	}
	u, _ := v.(*entity.CallUser)
	return u
}

// AddUserToCache 新建通话用户后写入缓存
func AddUserToCache(u *entity.CallUser) {
	if u == nil || userCacheMgr == nil {
		return
	}
	userCacheMgr.FlushCache(u.ID, u)
}

// FlushUserCache 通话用户变更后刷新缓存
func FlushUserCache(u *entity.CallUser) {
	AddUserToCache(u)
}

// GetUserFromCache 仅从内存缓存读取通话用户,未命中不访问数据库
func GetUserFromCache(userId uint64) *entity.CallUser {
	if userId == 0 || userCacheMgr == nil {
		return nil
	}
	v := userCacheMgr.GetFromCache(userId)
	if v == nil {
		return nil
	}
	u, _ := v.(*entity.CallUser)
	return u
}

// PreloadCallUsersToCache 启动时批量预热通话用户缓存(仅初始化调用一次)
func PreloadCallUsersToCache(userIds []uint64) {
	if len(userIds) == 0 || userCacheMgr == nil {
		return
	}
	ctx := gctx.New()
	users := make([]*entity.CallUser, 0, len(userIds))
	err := g.Model(string(entity.TbCallUser)).Ctx(ctx).
		WhereIn(string(db.IdName), userIds).
		Scan(&users)
	if err != nil {
		g.Log().Errorf(ctx, "preload call users failed: %v", err)
		return
	}
	for _, user := range users {
		AddUserToCache(user)
	}
}
