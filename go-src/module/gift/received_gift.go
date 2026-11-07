package gift

import (
	"context"

	"xr-game-server/constants/liverevenue"
	"xr-game-server/core/event"
	"xr-game-server/dao/anchorreceivedgiftdao"
	"xr-game-server/dto/giftdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
)

// InitReceivedGift 订阅礼物收益事件,记录主播收到过的礼物ID与数量
func InitReceivedGift() {
	event.Sub(gameevent.RevenueEventEvent, onRevenueEventForReceivedGift)
}

func onRevenueEventForReceivedGift(data any) {
	log, ok := data.(*entity.LiveRevenueLog)
	if !ok || log == nil {
		return
	}
	if log.Status != entity.LiveRevenueLogStatusNormal {
		return
	}
	rt := liverevenue.Type(log.RevenueType)
	if rt != liverevenue.Gift && rt != liverevenue.GiftLucky {
		return
	}
	if log.RoomId == 0 || log.BizId == 0 {
		return
	}
	addCount := uint64(log.Count)
	if addCount == 0 {
		addCount = 1
	}
	_ = anchorreceivedgiftdao.RecordReceivedGift(log.RoomId, log.BizId, addCount)
}

// GetReceivedGiftIds App 查询主播收到过的礼物ID与数量(读缓存)
func GetReceivedGiftIds(_ context.Context, req *giftdto.ReceivedGiftIdsReq) (*giftdto.ReceivedGiftIdsRes, error) {
	if req == nil || req.AnchorId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	rows := anchorreceivedgiftdao.GetReceivedGifts(req.AnchorId)
	list := make([]*giftdto.ReceivedGiftItem, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.GiftId == 0 {
			continue
		}
		count := row.Count
		if count == 0 {
			count = 1
		}
		list = append(list, &giftdto.ReceivedGiftItem{
			GiftId: row.GiftId,
			Count:  count,
		})
	}
	return &giftdto.ReceivedGiftIdsRes{List: list}, nil
}
