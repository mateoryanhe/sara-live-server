package calldao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
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
