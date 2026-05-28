package liveroomdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var liveRecordAudienceCacheMgr *cache.CacheMgr

func initLiveRecordAudienceDao() {
	liveRecordAudienceCacheMgr = cache.NewCacheMgr()
}

func getLiveRecordAudienceById(id string, liveRecordId, userId uint64) *entity.LiveRecordAudience {
	if liveRecordAudienceCacheMgr == nil {
		return nil
	}
	v := liveRecordAudienceCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var row *entity.LiveRecordAudience
		_ = g.Model(string(entity.TbLiveRecordAudience)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewLiveRecordAudience(liveRecordId, userId), nil
		}
		return row, nil
	})
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.LiveRecordAudience)
	return row
}

// TryRecordLiveRecordAudience 记录观众进入本场直播;已统计过返回 false
func TryRecordLiveRecordAudience(liveRecordId, userId uint64) bool {
	if liveRecordId == 0 || userId == 0 {
		return false
	}
	id := entity.BuildLiveRecordAudienceId(liveRecordId, userId)
	if getLiveRecordAudienceById(id, liveRecordId, userId) != nil {
		return false
	}
	return true
}
