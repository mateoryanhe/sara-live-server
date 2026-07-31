package shortvideodao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/constants/db"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var shortVideoStatCacheMgr *cache.CacheMgr

func initShortVideoStatDao() {
	shortVideoStatCacheMgr = cache.NewCacheMgr()
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
		if video := GetShortVideoById(videoId); video != nil {
			title = video.Title
			publishedAt = video.CreatedAt
		}
		return entity.NewShortVideoStat(videoId, title, publishedAt), nil
	})
	if cacheData == nil {
		return nil
	}
	stat, _ := cacheData.(*entity.ShortVideoStat)
	return stat
}

// ListStatPageByAuthorId 分页查询指定作者短视频统计数据(直接查库,按发布时间降序,不查总数)
func ListStatPageByAuthorId(authorId uint64, pageIndex, pageSize int) ([]*entity.ShortVideoStat, error) {
	list := make([]*entity.ShortVideoStat, 0)
	if authorId == 0 {
		return list, nil
	}
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	ctx := gctx.New()
	videoTable := string(entity.TbShortVideo)
	statTable := string(entity.TbShortVideoStat)
	err := g.Model(statTable+" s").Ctx(ctx).
		InnerJoin(videoTable+" v", "v."+string(db.IdName)+" = s."+string(db.IdName)).
		Where("v."+string(entity.ShortVideoAuthorId)+" = ?", authorId).
		Fields("s.*").
		OrderDesc("v." + string(db.CreatedAtName)).
		OrderDesc("s." + string(db.IdName)).
		Limit(pageSize).
		Offset((pageIndex - 1) * pageSize).
		Scan(&list)
	if err != nil {
		return list, err
	}
	return list, nil
}

func DeleteByVideoId(videoId uint64) error {
	if videoId == 0 {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbShortVideoStat)).WherePri(videoId).Delete()
	if err != nil {
		return err
	}
	if shortVideoStatCacheMgr != nil {
		_, _ = shortVideoStatCacheMgr.Cache.Remove(gctx.New(), videoId)
	}
	return nil
}
