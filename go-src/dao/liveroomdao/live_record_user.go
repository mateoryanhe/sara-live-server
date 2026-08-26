package liveroomdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

var liveRecordUserCacheMgr *cache.RowCache[*entity.LiveRecordUser]

func initLiveRecordUserDao() {
	liveRecordUserCacheMgr = cache.NewRowCache[*entity.LiveRecordUser]()
}

func getLiveRecordUserById(id string, liveRecordId, userId uint64) *entity.LiveRecordUser {
	if liveRecordUserCacheMgr == nil {
		return nil
	}
	v := liveRecordUserCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.LiveRecordUser, error) {
		var row *entity.LiveRecordUser
		_ = g.Model(string(entity.TbLiveRecordUser)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewLiveRecordUser(liveRecordId, userId), nil
		}
		return row, nil
	})
	return v
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

// TryRecordLiveRecordNewFollower 记录本场直播新加粉丝;本场已统计过返回 false
func TryRecordLiveRecordNewFollower(liveRecordId, userId uint64) bool {
	if liveRecordId == 0 || userId == 0 {
		return false
	}
	id := entity.BuildLiveRecordUserId(liveRecordId, userId)
	data := getLiveRecordUserById(id, liveRecordId, userId)
	if data == nil || !data.FollowerAt.IsZero() {
		return false
	}
	data.SetFollowerAt(time.Now())
	return true
}
