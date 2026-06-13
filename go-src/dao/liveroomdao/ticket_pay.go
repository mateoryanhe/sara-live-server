package liveroomdao

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/cache"
	"xr-game-server/entity"
)

var liveRoomTicketPayCacheMgr *cache.CacheMgr

func initLiveRoomTicketPayDao() {
	liveRoomTicketPayCacheMgr = cache.NewCacheMgr()
}

func getLiveRoomTicketPayById(id string, userId, roomId uint64) *entity.LiveRoomTicketPay {
	if liveRoomTicketPayCacheMgr == nil {
		return nil
	}
	v := liveRoomTicketPayCacheMgr.GetData(id, func(ctx context.Context) (value interface{}, err error) {
		var row *entity.LiveRoomTicketPay
		_ = g.Model(string(entity.TbLiveRoomTicketPay)).Where("id = ?", id).Scan(&row)
		if row == nil {
			return entity.NewLiveRoomTicketPay(userId, roomId), nil
		}
		return row, nil
	})
	if v == nil {
		return nil
	}
	row, _ := v.(*entity.LiveRoomTicketPay)
	return row
}

// GetLiveRoomTicketPay 获取用户在某直播间的门票扣费记录
func GetLiveRoomTicketPay(userId, roomId uint64) *entity.LiveRoomTicketPay {
	if userId == 0 || roomId == 0 {
		return nil
	}
	return getLiveRoomTicketPayById(entity.BuildLiveRoomTicketPayId(userId, roomId), userId, roomId)
}
