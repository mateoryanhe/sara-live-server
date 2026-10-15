package userwithdrawdao

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	userentity "xr-game-server/entity/user"
)

const userWithdrawLogListOrder = "updated_at desc, id desc"

const (
	// AppUserWithdrawLogListCachePageSize App 提现记录列表固定分页大小
	AppUserWithdrawLogListCachePageSize   = 20
	appUserWithdrawLogListCachedPages     = 2
	appUserWithdrawLogListCachedReadPages = appUserWithdrawLogListCachedPages
	appUserWithdrawLogListCacheMaxSize    = AppUserWithdrawLogListCachePageSize * appUserWithdrawLogListCachedPages
)

var appUserWithdrawLogListCacheMgr *cache.ListCache[*userentity.UserWithdrawLog]

func initAppUserWithdrawLogListCache() {
	appUserWithdrawLogListCacheMgr = cache.NewListCache[*userentity.UserWithdrawLog]()
}

func appUserWithdrawLogListCacheKey(userId uint64) string {
	return fmt.Sprintf("user_withdraw_app_list:%d", userId)
}

// ensureAppUserWithdrawLogListCacheLoaded 缓存未命中时按 updated_at 查库(前2页)并写入缓存
func ensureAppUserWithdrawLogListCacheLoaded(userId uint64) []*userentity.UserWithdrawLog {
	return getAppUserWithdrawLogListCache(userId)
}

func getAppUserWithdrawLogListCache(userId uint64) []*userentity.UserWithdrawLog {
	if appUserWithdrawLogListCacheMgr == nil || userId == 0 {
		return make([]*userentity.UserWithdrawLog, 0)
	}
	v := appUserWithdrawLogListCacheMgr.MustGetList(gctx.New(), appUserWithdrawLogListCacheKey(userId), func(ctx context.Context) ([]*userentity.UserWithdrawLog, error) {
		return loadAppUserWithdrawLogsFromDB(userId, 1, appUserWithdrawLogListCacheMaxSize), nil
	})
	if v == nil {
		return make([]*userentity.UserWithdrawLog, 0)
	}
	return v
}

func putAppUserWithdrawLogListCache(userId uint64, list []*userentity.UserWithdrawLog) {
	if appUserWithdrawLogListCacheMgr == nil || userId == 0 {
		return
	}
	if list == nil {
		list = make([]*userentity.UserWithdrawLog, 0)
	}
	appUserWithdrawLogListCacheMgr.PublishList(gctx.New(), appUserWithdrawLogListCacheKey(userId), list)
}

func cloneUserWithdrawLog(src *userentity.UserWithdrawLog) *userentity.UserWithdrawLog {
	if src == nil {
		return nil
	}
	cp := *src
	return &cp
}

// PrependUserWithdrawToAppListCache 按 userId 将记录置于 App 列表缓存头部(updated_at 最新在前)
func PrependUserWithdrawToAppListCache(row *userentity.UserWithdrawLog) {
	if appUserWithdrawLogListCacheMgr == nil || row == nil || row.UserId == 0 || row.ID == 0 {
		return
	}
	userId := row.UserId
	list := getAppUserWithdrawLogListCache(userId)
	newList := make([]*userentity.UserWithdrawLog, 0, len(list)+1)
	newList = append(newList, cloneUserWithdrawLog(row))
	for _, item := range list {
		if item != nil && item.ID != row.ID {
			newList = append(newList, item)
		}
	}
	if len(newList) > appUserWithdrawLogListCacheMaxSize {
		newList = newList[:appUserWithdrawLogListCacheMaxSize]
	}
	putAppUserWithdrawLogListCache(userId, newList)
}

// GetAppUserWithdrawLogsByUser 分页获取用户提现记录; 缓存2页, 前2页且 pageSize=20 时走缓存
func GetAppUserWithdrawLogsByUser(userId uint64, page, pageSize int) []*userentity.UserWithdrawLog {
	list := make([]*userentity.UserWithdrawLog, 0)
	if userId == 0 {
		return list
	}
	if page <= 0 {
		page = 1
	}
	pageSize = AppUserWithdrawLogListCachePageSize
	if page <= appUserWithdrawLogListCachedReadPages {
		cached := getAppUserWithdrawLogListCache(userId)
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
	return loadAppUserWithdrawLogsFromDB(userId, page, pageSize)
}

func loadAppUserWithdrawLogsFromDB(userId uint64, page, pageSize int) []*userentity.UserWithdrawLog {
	list := make([]*userentity.UserWithdrawLog, 0)
	if userId == 0 {
		return list
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = AppUserWithdrawLogListCachePageSize
	}
	_ = g.Model(string(userentity.TbUserWithdrawLog)).Ctx(gctx.New()).
		Where(string(userentity.UserWithdrawLogUserId)+" = ?", userId).
		Order(userWithdrawLogListOrder).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&list)
	return list
}

func CountUserWithdrawLogsByUser(userId uint64) (int, error) {
	if userId == 0 {
		return 0, nil
	}
	n, err := g.Model(string(userentity.TbUserWithdrawLog)).Ctx(gctx.New()).
		Where(string(userentity.UserWithdrawLogUserId)+" = ?", userId).
		Count()
	return n, err
}
