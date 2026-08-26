package liveroomdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity/live"
)

var liveRoomBillingPayCacheMgr *cache.RowCache[*entity.LiveRoomBillingPay]

func initLiveRoomBillingPayDao() {
	liveRoomBillingPayCacheMgr = cache.NewRowCache[*entity.LiveRoomBillingPay]()
}

func getLiveRoomBillingPayById(id string, userId, roomId uint64) *entity.LiveRoomBillingPay {
	if liveRoomBillingPayCacheMgr == nil {
		return nil
	}
	v := liveRoomBillingPayCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*entity.LiveRoomBillingPay, error) {
		var row *entity.LiveRoomBillingPay
		_ = g.Model(string(entity.TbLiveRoomBillingPay)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewLiveRoomBillingPay(userId, roomId), nil
		}
		return row, nil
	})
	return v
}

// GetLiveRoomBillingPay 获取观众在某场直播的按分钟计费记录
func GetLiveRoomBillingPay(userId, roomId uint64) *entity.LiveRoomBillingPay {
	if userId == 0 || roomId == 0 {
		return nil
	}
	return getLiveRoomBillingPayById(
		entity.BuildLiveRoomBillingPayId(userId, roomId),
		userId, roomId,
	)
}
