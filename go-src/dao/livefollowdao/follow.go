package livefollowdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/entity"
)

var (
	// followCacheMgr 按复合ID(userId_anchorId)缓存单条关注记录
	followCacheMgr *cache.CacheMgr
)

// InitLiveFollowDao 初始化关注主播相关缓存
func InitLiveFollowDao() {
	followCacheMgr = cache.NewCacheMgr()
}

// GetById 按复合ID获取(走缓存)
func GetById(id string) *entity.LiveFollow {
	v := followCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var f *entity.LiveFollow
		_ = g.Model(string(entity.TbLiveFollow)).Ctx(ctx).Where("id = ?", id).Scan(&f)
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

// GetFollowingsByUser 分页获取某用户当前已关注的记录(仅 Status == Follow)
func GetFollowingsByUser(userId uint64, page, pageSize int) (int, []*entity.LiveFollow) {
	list := make([]*entity.LiveFollow, 0)
	if userId == 0 {
		return 0, list
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	m := g.Model(string(entity.TbLiveFollow)).
		Where("user_id = ? AND status = ?", userId, entity.LiveFollowStatusFollow)
	total := userinfodao.GetFollowCount(userId)
	_ = m.Order("updated_at desc").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&list)
	return total, list
}

// GetFollowersByAnchor 分页获取某主播的粉丝记录(仅 Status == Follow)
func GetFollowersByAnchor(anchorId uint64, page, pageSize int) (int, []*entity.LiveFollow) {
	list := make([]*entity.LiveFollow, 0)
	if anchorId == 0 {
		return 0, list
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	m := g.Model(string(entity.TbLiveFollow)).
		Where("anchor_id = ? AND status = ?", anchorId, entity.LiveFollowStatusFollow)
	total := userinfodao.GetFollowerCount(anchorId)
	_ = m.Order("updated_at desc").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&list)
	return total, list
}

// CountFollowingsByUser 统计用户当前关注数(读 user_ext 缓存)
func CountFollowingsByUser(userId uint64) int {
	return userinfodao.GetFollowCount(userId)
}

// CountFollowersByAnchor 统计主播当前粉丝数(读 user_ext 缓存)
func CountFollowersByAnchor(anchorId uint64) int {
	return userinfodao.GetFollowerCount(anchorId)
}

// IsFollowing 查询 userId 是否已关注 anchorId
func IsFollowing(userId, anchorId uint64) bool {
	if userId == 0 || anchorId == 0 || userId == anchorId {
		return false
	}
	existing := GetByUserAnchor(userId, anchorId)
	return existing != nil && existing.Status == entity.LiveFollowStatusFollow
}

// AddFollowToCache 关注成功后刷新单条缓存
func AddFollowToCache(f *entity.LiveFollow) {
	if f == nil || followCacheMgr == nil {
		return
	}
	followCacheMgr.FlushCache(f.ID, f)
}
