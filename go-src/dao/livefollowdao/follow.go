package livefollowdao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var (
	// followCacheMgr 按复合ID(userId_anchorId)缓存单条关注记录
	followCacheMgr *cache.CacheMgr
	// userFollowsCacheMgr 按 userId 缓存"该用户关注的全部记录"列表
	userFollowsCacheMgr *cache.CacheMgr
	// anchorFollowersCacheMgr 按 anchorId 缓存"该主播的全部关注记录"列表
	anchorFollowersCacheMgr *cache.CacheMgr
)

// InitLiveFollowDao 初始化关注主播相关缓存
func InitLiveFollowDao() {
	followCacheMgr = cache.NewCacheMgr()
	userFollowsCacheMgr = cache.NewCacheMgr()
	anchorFollowersCacheMgr = cache.NewCacheMgr()
}

// GetById 按复合ID获取(走缓存)
func GetById(id string) *entity.LiveFollow {
	v := followCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var f *entity.LiveFollow
		_ = g.Model(string(entity.TbLiveFollow)).Where("id = ?", id).Scan(&f)
		return f, nil
	})
	if v == nil {
		return nil
	}
	f, _ := v.(*entity.LiveFollow)
	return f
}

// GetByUserAnchor 按 (userId, anchorId) 获取(走缓存)
func GetByUserAnchor(userId, anchorId uint64) *entity.LiveFollow {
	return GetById(entity.BuildLiveFollowId(userId, anchorId))
}

// GetFollowingsByUser 获取某用户当前已关注的全部记录(仅 Status == Follow)
func GetFollowingsByUser(userId uint64) []*entity.LiveFollow {
	v := userFollowsCacheMgr.GetData(userId, func(ctx context.Context) (value interface{}, err error) {
		list := make([]*entity.LiveFollow, 0)
		_ = g.Model(string(entity.TbLiveFollow)).
			Where("user_id = ? AND status = ?", userId, entity.LiveFollowStatusFollow).
			Order("updated_at desc").
			Scan(&list)
		return list, nil
	})
	if v == nil {
		return make([]*entity.LiveFollow, 0)
	}
	list, _ := v.([]*entity.LiveFollow)
	if list == nil {
		return make([]*entity.LiveFollow, 0)
	}
	return list
}

// GetFollowersByAnchor 获取某主播的全部粉丝记录(仅 Status == Follow)
func GetFollowersByAnchor(anchorId uint64) []*entity.LiveFollow {
	v := anchorFollowersCacheMgr.GetData(anchorId, func(ctx context.Context) (value interface{}, err error) {
		list := make([]*entity.LiveFollow, 0)
		_ = g.Model(string(entity.TbLiveFollow)).
			Where("anchor_id = ? AND status = ?", anchorId, entity.LiveFollowStatusFollow).
			Order("updated_at desc").
			Scan(&list)
		return list, nil
	})
	if v == nil {
		return make([]*entity.LiveFollow, 0)
	}
	list, _ := v.([]*entity.LiveFollow)
	if list == nil {
		return make([]*entity.LiveFollow, 0)
	}
	return list
}

// AddFollowToCache 关注成功后写入(单条 + 双向列表)
func AddFollowToCache(f *entity.LiveFollow) {
	if f == nil {
		return
	}
	followCacheMgr.FlushCache(f.ID, f)

	// 用户关注列表
	userList := GetFollowingsByUser(f.UserId)
	exists := false
	for _, item := range userList {
		if item.ID == f.ID {
			exists = true
			break
		}
	}
	if !exists {
		userList = append([]*entity.LiveFollow{f}, userList...)
		userFollowsCacheMgr.FlushCache(f.UserId, userList)
	}

	// 主播粉丝列表
	anchorList := GetFollowersByAnchor(f.AnchorId)
	exists = false
	for _, item := range anchorList {
		if item.ID == f.ID {
			exists = true
			break
		}
	}
	if !exists {
		anchorList = append([]*entity.LiveFollow{f}, anchorList...)
		anchorFollowersCacheMgr.FlushCache(f.AnchorId, anchorList)
	}
}

// RemoveFollowFromCache 取消关注后,从双向列表中剔除(单条缓存保留以便复用,Status 已切回 Unfollow)
func RemoveFollowFromCache(f *entity.LiveFollow) {
	if f == nil {
		return
	}
	followCacheMgr.FlushCache(f.ID, f)

	if userFollowsCacheMgr != nil {
		userList := GetFollowingsByUser(f.UserId)
		newList := make([]*entity.LiveFollow, 0, len(userList))
		for _, item := range userList {
			if item.ID == f.ID {
				continue
			}
			newList = append(newList, item)
		}
		userFollowsCacheMgr.FlushCache(f.UserId, newList)
	}

	if anchorFollowersCacheMgr != nil {
		anchorList := GetFollowersByAnchor(f.AnchorId)
		newList := make([]*entity.LiveFollow, 0, len(anchorList))
		for _, item := range anchorList {
			if item.ID == f.ID {
				continue
			}
			newList = append(newList, item)
		}
		anchorFollowersCacheMgr.FlushCache(f.AnchorId, newList)
	}
}
