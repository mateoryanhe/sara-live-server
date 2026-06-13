package liveroomdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var liveRoomBillingPayCacheMgr *cache.CacheMgr

func initLiveRoomBillingPayDao() {
	liveRoomBillingPayCacheMgr = cache.NewCacheMgr()
}

func getLiveRoomBillingPayById(id string, userId, roomId, liveRecordId uint64) *entity.LiveRoomBillingPay {
	if liveRoomBillingPayCacheMgr == nil {
		return nil
	}
	v := liveRoomBillingPayCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var row *entity.LiveRoomBillingPay
		_ = g.Model(string(entity.TbLiveRoomBillingPay)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewLiveRoomBillingPay(userId, roomId, liveRecordId), nil
		}
		return row, nil
	})
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.LiveRoomBillingPay)
	return row
}

// GetLiveRoomBillingPay 获取观众在某场直播的按分钟计费记录
func GetLiveRoomBillingPay(userId, roomId, liveRecordId uint64) *entity.LiveRoomBillingPay {
	if userId == 0 || roomId == 0 || liveRecordId == 0 {
		return nil
	}
	return getLiveRoomBillingPayById(
		entity.BuildLiveRoomBillingPayId(userId, roomId, liveRecordId),
		userId, roomId, liveRecordId,
	)
}
