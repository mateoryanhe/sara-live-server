package messagedao

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	"xr-game-server/core/str"
	"xr-game-server/dto/activitymessagedto"
	"xr-game-server/entity"
)

const activityMessageCacheKey = "all"

var activityMessageCacheMgr *cache.CacheMgr

func initActivityMessageDao() {
	activityMessageCacheMgr = cache.NewCacheMgr()
}

func loadAllActivityMessagesFromDB() []*entity.ActivityMessage {
	ret := make([]*entity.ActivityMessage, 0)
	_ = g.DB().Model(string(entity.TbActivityMessage)).
		Order("published_at desc, created_at desc").
		Scan(&ret)
	return ret
}

// GetAllActivityMessagesCached 加载全部活动消息(优先 gcache,未命中再查库)
func GetAllActivityMessagesCached() []*entity.ActivityMessage {
	if activityMessageCacheMgr == nil {
		return loadAllActivityMessagesFromDB()
	}
	v := activityMessageCacheMgr.GetData(activityMessageCacheKey, func(ctx context.Context) (interface{}, error) {
		return loadAllActivityMessagesFromDB(), nil
	})
	if v == nil {
		return nil
	}
	rows, _ := v.([]*entity.ActivityMessage)
	return rows
}

// RemoveActivityMessageCache 活动消息变更后移除 gcache
func RemoveActivityMessageCache() {
	if activityMessageCacheMgr == nil {
		return
	}
	_, _ = activityMessageCacheMgr.Cache.Remove(gctx.New(), activityMessageCacheKey)
}

func GetActivityMessageById(id uint64) *entity.ActivityMessage {
	var row entity.ActivityMessage
	if err := g.DB().Model(string(entity.TbActivityMessage)).Where("id = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func CreateActivityMessage(row *entity.ActivityMessage) error {
	_, err := g.DB().Model(string(entity.TbActivityMessage)).Save(row)
	if err == nil {
		RemoveActivityMessageCache()
	}
	return err
}

func UpdateActivityMessage(row *entity.ActivityMessage) error {
	return CreateActivityMessage(row)
}

func DeleteActivityMessage(id uint64) error {
	_, err := g.DB().Model(string(entity.TbActivityMessage)).WherePri(id).Delete()
	if err == nil {
		RemoveActivityMessageCache()
	}
	return err
}

func GetActivityMessagesByIDs(ids []uint64) []*entity.ActivityMessage {
	if len(ids) == 0 {
		return nil
	}
	var rows []*entity.ActivityMessage
	_ = g.DB().Model(string(entity.TbActivityMessage)).WhereIn("id", ids).Scan(&rows)
	return rows
}

// ReloadActivityMessageCaches 活动消息数据变更后刷新相关内存缓存
func ReloadActivityMessageCaches() {
	RemoveActivityMessageCache()
	ClearAllUserActivityMessageListCacheA()
}

func GetActivityMessageList(req *activitymessagedto.ActivityMessageListReq) (int, []*activitymessagedto.ActivityMessageListRes) {
	sql := `select id,
                   icon_en, icon_es, icon_pt, icon_hi,
                   bg_en, bg_es, bg_pt, bg_hi,
                   title_en, title_es, title_pt, title_hi,
                   content_en, content_es, content_pt, content_hi,
                   url_en, url_es, url_pt, url_hi,
                   status, published_at, created_at, updated_at
            from activity_messages
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*activitymessagedto.ActivityMessageListRes, 0)

	if req.Title != "" {
		sql += ` and (title_en LIKE ? or title_es LIKE ? or title_pt LIKE ? or title_hi LIKE ?)`
		like := fmt.Sprintf("%%%s%%", req.Title)
		param = append(param, like, like, like, like)
	}
	switch req.StatusFilter {
	case 1:
		sql += ` and status = ?`
		param = append(param, entity.ActivityMessageStatusUnpublished)
	case 2:
		sql += ` and status = ?`
		param = append(param, entity.ActivityMessageStatusPublished)
	}

	sql += ` order by published_at desc, created_at desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
