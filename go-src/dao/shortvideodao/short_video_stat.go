package shortvideodao

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

const (
	StatListCachePageSize   = 50
	statListCacheMaxPages   = 5 // n: 缓存 n 页数据
	statListCachedReadPages = statListCacheMaxPages - 1
	statListCacheMaxSize    = StatListCachePageSize * statListCacheMaxPages
)

var (
	shortVideoStatCacheMgr *cache.CacheMgr
	authorStatListCacheMgr *cache.CacheMgr
)

func initShortVideoStatDao() {
	shortVideoStatCacheMgr = cache.NewPermanentCacheMgr()
	authorStatListCacheMgr = cache.NewCacheMgr()
}

func authorStatListCacheKey(authorId uint64) string {
	return fmt.Sprintf("short_video_stat_list:%d", authorId)
}

// GetStatByVideoId 根据视频ID获取统计数据,不存在则新建内存对象并写入标题
func GetStatByVideoId(videoId uint64) *entity.ShortVideoStat {
	if videoId == 0 || shortVideoStatCacheMgr == nil {
		return nil
	}
	cacheData := shortVideoStatCacheMgr.GetData(videoId, func(ctx context.Context) (value interface{}, err error) {
		var row *entity.ShortVideoStat
		err = g.Model(string(entity.TbShortVideoStat)).Where(g.Map{
			string(db.IdName): videoId,
		}).Scan(&row)
		if row != nil && row.ID != 0 {
			return row, nil
		}
		title := ""
		publishedAt := time.Now()
		authorId := uint64(0)
		if video := GetShortVideoById(videoId); video != nil {
			title = video.Title
			publishedAt = video.CreatedAt
			authorId = video.AuthorId
		}
		return entity.NewShortVideoStat(videoId, authorId, title, publishedAt), nil
	})
	if cacheData == nil {
		return nil
	}

	stat, _ := cacheData.(*entity.ShortVideoStat)
	return stat
}

// GetStatFromCacheByVideoId 仅从内存缓存读取统计数据,未命中不访问数据库
func GetStatFromCacheByVideoId(videoId uint64) *entity.ShortVideoStat {
	if videoId == 0 || shortVideoStatCacheMgr == nil {
		return nil
	}
	v := shortVideoStatCacheMgr.GetFromCache(videoId)
	if v == nil {
		return nil
	}
	stat, _ := v.(*entity.ShortVideoStat)
	return stat
}

func getAuthorStatListCache(authorId uint64) []*entity.ShortVideoStat {
	if authorStatListCacheMgr == nil || authorId == 0 {
		return make([]*entity.ShortVideoStat, 0)
	}
	v := authorStatListCacheMgr.GetData(authorStatListCacheKey(authorId), func(ctx context.Context) (interface{}, error) {
		return loadStatListPageFromDB(authorId, 1, statListCacheMaxSize), nil
	})
	list, _ := v.([]*entity.ShortVideoStat)
	if list == nil {
		return make([]*entity.ShortVideoStat, 0)
	}
	return list
}

func putAuthorStatListCache(authorId uint64, list []*entity.ShortVideoStat) {
	if authorStatListCacheMgr == nil || authorId == 0 {
		return
	}
	if list == nil {
		list = make([]*entity.ShortVideoStat, 0)
	}
	authorStatListCacheMgr.FlushCache(authorStatListCacheKey(authorId), list)
}

// PrependStatToAuthorListCache CMS审核通过时将统计记录写入作者列表缓存并按ID降序排序
func PrependStatToAuthorListCache(authorId uint64, stat *entity.ShortVideoStat) {
	if authorStatListCacheMgr == nil || authorId == 0 || stat == nil || stat.ID == 0 {
		return
	}
	list := getAuthorStatListCache(authorId)
	newList := make([]*entity.ShortVideoStat, 0, len(list)+1)
	newList = append(newList, stat)
	for _, row := range list {
		if row != nil && row.ID != stat.ID {
			newList = append(newList, row)
		}
	}
	sortStatListByIdDesc(newList)
	if len(newList) > statListCacheMaxSize {
		newList = newList[:statListCacheMaxSize]
	}
	putAuthorStatListCache(authorId, newList)
}

func sortStatListByIdDesc(list []*entity.ShortVideoStat) {
	sort.Slice(list, func(i, j int) bool {
		a, b := list[i], list[j]
		if a == nil || b == nil {
			return a != nil
		}
		return a.ID > b.ID
	})
}

func loadStatListPageFromDB(authorId uint64, pageIndex, pageSize int) []*entity.ShortVideoStat {
	list := make([]*entity.ShortVideoStat, 0)
	if authorId == 0 {
		return list
	}
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = StatListCachePageSize
	}
	_ = g.Model(string(entity.TbShortVideoStat)).
		Where(string(entity.ShortVideoStatAuthorId)+" = ?", authorId).
		OrderDesc(string(db.IdName)).
		Limit(pageSize).
		Offset((pageIndex - 1) * pageSize).
		Scan(&list)
	return list
}

func sliceStatListPage(list []*entity.ShortVideoStat, pageIndex, pageSize int) []*entity.ShortVideoStat {
	if len(list) == 0 {
		return make([]*entity.ShortVideoStat, 0)
	}
	start := (pageIndex - 1) * pageSize
	if start >= len(list) {
		return make([]*entity.ShortVideoStat, 0)
	}
	end := start + pageSize
	if end > len(list) {
		end = len(list)
	}
	return list[start:end]
}

// ListStatPageByAuthorId 分页查询指定作者短视频统计数据(单表查询,按ID降序,不查总数)
// 缓存 n 页数据,前 n-1 页且 pageSize=StatListCachePageSize 时走缓存,第 n 页起直接查库
func ListStatPageByAuthorId(authorId uint64, pageIndex, pageSize int) ([]*entity.ShortVideoStat, error) {
	list := make([]*entity.ShortVideoStat, 0)
	if authorId == 0 {
		return list, nil
	}
	if pageIndex <= 0 {
		pageIndex = 1
	}
	pageSize = StatListCachePageSize
	if pageIndex <= statListCachedReadPages {
		return sliceStatListPage(getAuthorStatListCache(authorId), pageIndex, pageSize), nil
	}
	return loadStatListPageFromDB(authorId, pageIndex, pageSize), nil
}

// RemoveStatCacheByVideoId CMS删除视频时移除单条统计缓存
func RemoveStatCacheByVideoId(videoId uint64) {
	if shortVideoStatCacheMgr == nil || videoId == 0 {
		return
	}
	_, _ = shortVideoStatCacheMgr.Cache.Remove(gctx.New(), videoId)
}

// RefreshAuthorStatListCacheOnVideoDelete CMS删除视频时刷新作者统计列表缓存(移除该视频,其余项用最新统计替换并排序)
func RefreshAuthorStatListCacheOnVideoDelete(authorId, videoId uint64) {
	if authorStatListCacheMgr == nil || authorId == 0 || videoId == 0 {
		return
	}
	key := authorStatListCacheKey(authorId)
	v := authorStatListCacheMgr.GetFromCache(key)
	if v == nil {
		return
	}
	list, _ := v.([]*entity.ShortVideoStat)
	if len(list) == 0 {
		return
	}
	found := false
	newObj := GetStatByVideoId(videoId)
	newList := make([]*entity.ShortVideoStat, 0, len(list))
	for _, row := range list {
		if row == nil {
			continue
		}
		if row.ID == videoId {
			found = true
			newList = append(newList, newObj)
		} else {
			newList = append(newList, row)
		}

	}
	if !found {
		return
	}
	putAuthorStatListCache(authorId, newList)
}
