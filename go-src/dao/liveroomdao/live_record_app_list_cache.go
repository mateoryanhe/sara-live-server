package liveroomdao

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

const (
	// AppLiveRecordListCachePageSize App 端直播记录列表固定分页大小
	AppLiveRecordListCachePageSize   = 50
	appLiveRecordListCachedPages     = 3
	appLiveRecordListCachedReadPages = appLiveRecordListCachedPages - 1
	appLiveRecordListCacheMaxSize    = AppLiveRecordListCachePageSize * appLiveRecordListCachedPages
)

var appLiveRecordListCacheMgr *cache.CacheMgr

func initAppLiveRecordListCache() {
	appLiveRecordListCacheMgr = cache.NewCacheMgr()
}

func appLiveRecordListCacheKey(anchorId uint64) string {
	return fmt.Sprintf("live_record_app_list:%d", anchorId)
}

func getAppLiveRecordListCache(anchorId uint64) []*entity.LiveRecord {
	if appLiveRecordListCacheMgr == nil || anchorId == 0 {
		return make([]*entity.LiveRecord, 0)
	}
	v := appLiveRecordListCacheMgr.GetData(appLiveRecordListCacheKey(anchorId), func(ctx context.Context) (interface{}, error) {
		return loadAppLiveRecordsFromDB(anchorId, 1, appLiveRecordListCacheMaxSize), nil
	})
	list, _ := v.([]*entity.LiveRecord)
	if list == nil {
		return make([]*entity.LiveRecord, 0)
	}
	return list
}

func putAppLiveRecordListCache(anchorId uint64, list []*entity.LiveRecord) {
	if appLiveRecordListCacheMgr == nil || anchorId == 0 {
		return
	}
	if list == nil {
		list = make([]*entity.LiveRecord, 0)
	}
	appLiveRecordListCacheMgr.FlushCache(appLiveRecordListCacheKey(anchorId), list)
}

// PrependLiveRecordToAppListCache 新开直播后将记录插入列表缓存头部
func PrependLiveRecordToAppListCache(anchorId uint64, record *entity.LiveRecord) {
	if appLiveRecordListCacheMgr == nil || anchorId == 0 || record == nil || record.ID == 0 {
		return
	}
	list := getAppLiveRecordListCache(anchorId)
	newList := make([]*entity.LiveRecord, 0, len(list)+1)
	newList = append(newList, record)
	for _, row := range list {
		if row != nil && row.ID != record.ID {
			newList = append(newList, row)
		}
	}
	if len(newList) > appLiveRecordListCacheMaxSize {
		newList = newList[:appLiveRecordListCacheMaxSize]
	}
	putAppLiveRecordListCache(anchorId, newList)
}

// GetAppLiveRecordsByAnchor 分页获取主播直播记录
// 缓存3页(150条),前2页且 pageSize=AppLiveRecordListCachePageSize 时走缓存,第3页起直接查库
func GetAppLiveRecordsByAnchor(anchorId uint64, page, pageSize int) []*entity.LiveRecord {
	list := make([]*entity.LiveRecord, 0)
	if anchorId == 0 {
		return list
	}
	if page <= 0 {
		page = 1
	}
	pageSize = AppLiveRecordListCachePageSize
	if page <= appLiveRecordListCachedReadPages {
		cached := getAppLiveRecordListCache(anchorId)
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
	return loadAppLiveRecordsFromDB(anchorId, page, pageSize)
}

func loadAppLiveRecordsFromDB(anchorId uint64, page, pageSize int) []*entity.LiveRecord {
	list := make([]*entity.LiveRecord, 0)
	if anchorId == 0 {
		return list
	}
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = AppLiveRecordListCachePageSize
	}
	_ = g.Model(string(entity.TbLiveRecord)).Ctx(gctx.New()).
		Where(string(entity.LiveRecordAnchorId)+" = ?", anchorId).
		Order("id desc").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&list)
	return list
}
