package shortvideodao

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var watchOneCacheMgr *cache.CacheMgr

func initShortVideoWatchDao() {
	watchOneCacheMgr = cache.NewCacheMgr()
}

func watchCacheKey(userId, videoId uint64) string {
	return fmt.Sprintf("%v_%v", userId, videoId)
}

func GetOneShortVideoWatch(userId, videoId uint64) *entity.ShortVideoWatch {
	id := watchCacheKey(userId, videoId)
	v := watchOneCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var watch entity.ShortVideoWatch
		if scanErr := g.DB().Model(string(entity.TbShortVideoWatch)).
			Where("user_id = ? AND video_id = ?", userId, videoId).
			Scan(&watch); scanErr == nil && watch.ID != "" {
			return &watch, nil
		}
		return entity.NewShortVideoWatch(userId, videoId), nil
	})
	return v.(*entity.ShortVideoWatch)
}

// GetShortVideoWatch 分页查询用户可继续观看的记录(已上架且免费/已付费/作者本人/仍有试看)
func GetShortVideoWatch(userId uint64, pageIndex, pageSize int) []*entity.ShortVideoWatch {
	if userId == 0 {
		return make([]*entity.ShortVideoWatch, 0)
	}
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	ctx := gctx.New()
	whereSQL := `FROM short_video_watches w
INNER JOIN short_videos v ON v.id = w.video_id
WHERE w.user_id = ?
AND v.status = ?
AND (
	v.is_paid != ?
	OR v.author_id = ?
	OR w.paid_time IS NOT NULL
	
)`
	args := []any{
		userId,
		entity.ShortVideoStatusOnShelf,
		entity.ShortVideoPaidYes,
		userId,
	}

	list := make([]*entity.ShortVideoWatch, 0)
	querySQL := `SELECT w.* ` + whereSQL + `
ORDER BY w.updated_at DESC, w.id DESC
LIMIT ? OFFSET ?`
	queryArgs := append(append([]any{}, args...), pageSize, (pageIndex-1)*pageSize)
	_ = g.DB().GetScan(ctx, &list, querySQL, queryArgs...)
	return list
}

// DeleteWatchByVideoId 删除指定视频的观看/点赞记录,并清理单条缓存
func DeleteWatchByVideoId(videoId uint64) error {
	if videoId == 0 {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbShortVideoWatch)).Where("video_id = ?", videoId).Delete()
	if err != nil {
		return err
	}
	removeWatchCacheByVideoId(videoId)
	return nil
}

func removeWatchCacheByVideoId(videoId uint64) {
	if watchOneCacheMgr == nil || videoId == 0 {
		return
	}
	ctx := gctx.New()
	keys, err := watchOneCacheMgr.Cache.Keys(ctx)
	if err != nil || len(keys) == 0 {
		return
	}
	videoIdStr := strconv.FormatUint(videoId, 10)
	for _, key := range keys {
		keyStr, ok := key.(string)
		if !ok {
			continue
		}
		idx := strings.LastIndex(keyStr, "_")
		if idx < 0 || keyStr[idx+1:] != videoIdStr {
			continue
		}
		_, _ = watchOneCacheMgr.Cache.Remove(ctx, key)
	}
}
