package guilddao

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var guildCacheMgr *cache.CacheMgr

func InitGuildDao() {
	guildCacheMgr = cache.NewCacheMgr()
}

// GetGuildByIdCached 走缓存读取工会信息(给app使用)
func GetGuildByIdCached(id uint64) *entity.LiveGuild {
	v := guildCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var guild *entity.LiveGuild
		_ = g.Model(string(entity.TbLiveGuild)).Where("id = ?", id).Scan(&guild)
		return guild, nil
	})
	if v == nil {
		return nil
	}
	guild, _ := v.(*entity.LiveGuild)
	return guild
}

// RemoveGuildCache 移除工会缓存
func RemoveGuildCache(id uint64) {
	if guildCacheMgr == nil {
		return
	}
	_, _ = guildCacheMgr.Cache.Remove(gctx.New(), id)
}
