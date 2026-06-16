package liveroomdao

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var liveRecordUserCacheMgr *cache.CacheMgr

func initLiveRecordUserDao() {
	liveRecordUserCacheMgr = cache.NewCacheMgr()
}

func getLiveRecordUserById(id string, liveRecordId, userId uint64) *entity.LiveRecordUser {
	if liveRecordUserCacheMgr == nil {
		return nil
	}
	v := liveRecordUserCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var row *entity.LiveRecordUser
		_ = g.Model(string(entity.TbLiveRecordUser)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewLiveRecordUser(liveRecordId, userId), nil
		}
		return row, nil
	})
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.LiveRecordUser)
	return row
}

// TryRecordLiveRecordAudience 记录观众进入本场直播;本场已统计过返回 false
func TryRecordLiveRecordAudience(liveRecordId, userId uint64) bool {
	if liveRecordId == 0 || userId == 0 {
		return false
	}
	id := entity.BuildLiveRecordUserId(liveRecordId, userId)
	data := getLiveRecordUserById(id, liveRecordId, userId)
	if data == nil || !data.AudienceAt.IsZero() {
		return false
	}
	data.SetAudienceAt(time.Now())
	return true
}

// TryRecordLiveRecordGiftSender 记录本场直播送礼人;本场已统计过返回 false
func TryRecordLiveRecordGiftSender(liveRecordId, userId uint64) bool {
	if liveRecordId == 0 || userId == 0 {
		return false
	}
	id := entity.BuildLiveRecordUserId(liveRecordId, userId)
	data := getLiveRecordUserById(id, liveRecordId, userId)
	if data == nil || !data.GiftSenderAt.IsZero() {
		return false
	}
	data.SetGiftSenderAt(time.Now())
	return true
}
