package gamebetdao

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/game"
)

const (
	// AppGameBetListCachePageSize App 端游戏下注列表固定分页大小
	AppGameBetListCachePageSize   = 50
	appGameBetListCachedPages     = 3
	appGameBetListCachedReadPages = appGameBetListCachedPages
	appGameBetListCacheMaxSize    = AppGameBetListCachePageSize * appGameBetListCachedPages
)

var appGameBetListCacheMgr *cache.CacheMgr

func initAppGameBetListCache() {
	appGameBetListCacheMgr = cache.NewCacheMgr()
}

func appGameBetListCacheKey(userId uint64) string {
	return fmt.Sprintf("game_bet_app_list:%d", userId)
}

func getAppGameBetListCache(userId uint64) []*entity.GameBetLog {
	if appGameBetListCacheMgr == nil || userId == 0 {
		return make([]*entity.GameBetLog, 0)
	}
	v := appGameBetListCacheMgr.GetData(appGameBetListCacheKey(userId), func(ctx context.Context) (interface{}, error) {
		return loadAppGameBetLogsFromDB(userId, 1, appGameBetListCacheMaxSize), nil
	})
	list, _ := v.([]*entity.GameBetLog)
	if list == nil {
		return make([]*entity.GameBetLog, 0)
	}
	return list
}

func putAppGameBetListCache(userId uint64, list []*entity.GameBetLog) {
	if appGameBetListCacheMgr == nil || userId == 0 {
		return
	}
	if list == nil {
		list = make([]*entity.GameBetLog, 0)
	}
	appGameBetListCacheMgr.FlushCache(appGameBetListCacheKey(userId), list)
}

// PrependGameBetToAppListCache 新下注记录插入列表缓存头部,超出容量则截断
func PrependGameBetToAppListCache(userId uint64, row *entity.GameBetLog) {
	if appGameBetListCacheMgr == nil || userId == 0 || row == nil || row.ID == 0 {
		return
	}
	list := getAppGameBetListCache(userId)
	newList := make([]*entity.GameBetLog, 0, len(list)+1)
	newList = append(newList, row)
	for _, item := range list {
		if item != nil && item.ID != row.ID {
			newList = append(newList, item)
		}
	}
	if len(newList) > appGameBetListCacheMaxSize {
		newList = newList[:appGameBetListCacheMaxSize]
	}
	putAppGameBetListCache(userId, newList)
}

// GetAppGameBetLogsByUser 分页获取用户游戏下注记录
// 缓存3页(150条),前3页且 pageSize=50 时走缓存,第4页起直接查库
func GetAppGameBetLogsByUser(userId uint64, page, pageSize int) []*entity.GameBetLog {
	list := make([]*entity.GameBetLog, 0)
	if userId == 0 {
		return list
	}
	if page <= 0 {
		page = 1
	}
	pageSize = AppGameBetListCachePageSize
	if page <= appGameBetListCachedReadPages {
		cached := getAppGameBetListCache(userId)
		start := (page - 1) * pageSize
		if start >= len(cached) {
			return list
		}
		end := start + pageSize
		if end > len(cached) {
			end = len(cached)
		}
		return cached[start:end]
	}
	return loadAppGameBetLogsFromDB(userId, page, pageSize)
}

func loadAppGameBetLogsFromDB(userId uint64, page, pageSize int) []*entity.GameBetLog {
	list := make([]*entity.GameBetLog, 0)
	if userId == 0 {
		return list
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = AppGameBetListCachePageSize
	}
	_ = g.Model(string(entity.TbGameBetLog)).Ctx(gctx.New()).
		Where(string(entity.GameBetLogUserId)+" = ?", userId).
		Where(string(entity.GameBetLogAmount)+" > ?", 0).
		Order("id desc").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&list)
	return list
}
