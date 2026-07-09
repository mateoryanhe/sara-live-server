package userchanneltokendao

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var tokenCacheMgr *cache.CacheMgr

// Init 初始化用户频道Token DAO 缓存
func Init() {
	tokenCacheMgr = cache.NewCacheMgr()
}

func buildCacheKey(userId uint64, channelName string) string {
	return fmt.Sprintf("%d_%s", userId, channelName)
}

// GetByUserChannel 按 userId + channelName 查询(走缓存)
func GetByUserChannel(userId uint64, channelName string) *entity.UserChannelToken {
	if userId == 0 || channelName == "" || tokenCacheMgr == nil {
		return nil
	}
	cacheKey := buildCacheKey(userId, channelName)
	v := tokenCacheMgr.GetData(cacheKey, func(ctx context.Context) (value interface{}, err error) {
		var ret entity.UserChannelToken
		err = g.DB().Model(string(entity.TbUserChannelToken)).
			Where(string(entity.UserChannelTokenUserId)+" = ?", userId).
			Where(string(entity.UserChannelTokenChannelName)+" = ?", channelName).
			Scan(&ret)
		if err != nil || ret.UserId == 0 {
			return nil, err
		}
		return &ret, nil
	})
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.UserChannelToken)
	return row
}

// AddToCache 写入或刷新缓存
func AddToCache(row *entity.UserChannelToken) {
	if row == nil || tokenCacheMgr == nil || row.UserId == 0 || row.ChannelName == "" {
		return
	}
	tokenCacheMgr.FlushCache(buildCacheKey(row.UserId, row.ChannelName), row)
}

// FlushCache 记录变更后刷新缓存; row为nil时表示声网配置变更,清空全部Token
func FlushCache(row *entity.UserChannelToken) {
	if row == nil {
		flushAllOnAgoraCfgChanged()
		return
	}
	AddToCache(row)
}

func flushAllOnAgoraCfgChanged() {
	if tokenCacheMgr == nil {
		return
	}
	_, _ = g.DB().Model(string(entity.TbUserChannelToken)).Delete()
	refreshTokenCacheAfterDelete()
}

func refreshTokenCacheAfterDelete() {
	ctx := gctx.New()
	tokenCacheMgr.Cache.Clear(ctx)
}
