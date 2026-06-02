package liveroom

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gmlock"
	"time"
	"xr-game-server/constants/cmd"
	"xr-game-server/constants/currency"
	"xr-game-server/constants/liverevenue"
	"xr-game-server/core/event"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/push"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
	"xr-game-server/module/gift"
	"xr-game-server/module/upload"
	"xr-game-server/module/wallet"
)

// SendGift 直播间送礼
//  1. 校验房间存在、礼物存在(命中礼物缓存,即默认已上架)、数量合法
//  2. 计算总消耗 = 礼物单价 * 数量,使用钻石支付(diamond.Sub)
//  3. 扣款成功后,向房间内全体在线用户(含送礼人自身)推送 cmd.LiveRoomGift
func SendGift(ctx context.Context, req *liveroomdto.SendGiftReq) (*liveroomdto.SendGiftRes, error) {
	if req.Count <= 0 {
		return nil, errercode.CreateCode(errercode.GiftCountInvalid)
	}

	senderId := httpserver.GetAuthId(ctx)

	// 1. 房间存在性校验
	room := liveroomdao.GetRoomById(req.RoomId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	// 2. 礼物配置(从缓存读取,缓存仅包含已上架礼物)
	giftItem := gift.GetGiftFromCacheById(req.GiftId)
	if giftItem == nil {
		return nil, errercode.CreateCode(errercode.GiftOffShelf)
	}

	// 3. 计算总价并扣减钻石(uint64 防溢出)
	if giftItem.Price > 0 && uint64(req.Count) > (^uint64(0))/giftItem.Price {
		return nil, errercode.CreateCode(errercode.GiftCountInvalid)
	}
	totalCost := giftItem.Price * uint64(req.Count)
	remaining, err := wallet.DiamondSub(senderId, float64(totalCost), currency.ReasonGiftSend)
	if err != nil {
		return nil, err
	}

	//记录直播收益流水(礼物)
	eventData := entity.NewLiveRevenueLogRecord(room.ID, room.LiveRecordId, senderId, room.ID, req.GiftId, req.Count, giftItem.Price, totalCost, uint8(liverevenue.Gift))

	// 4. 构造推送载荷,广播给房间内所有在线用户
	sender := userinfodao.GetUserInfoByUserId(senderId)
	payload := &liveroomdto.GiftPushItem{
		RoomId:    req.RoomId,
		SenderId:  senderId,
		GiftId:    giftItem.ID,
		GiftName:  giftItem.Name,
		GiftIcon:  giftItem.Icon,
		GiftAnim:  giftItem.Animation,
		UnitPrice: giftItem.Price,
		Count:     req.Count,
		TotalCost: totalCost,
		SentAt:    time.Now().Unix(),
	}
	if sender != nil {
		payload.SenderName = sender.Nickname
		payload.SenderAvatar = upload.GetUrlByName(sender.Avatar)
	}

	for _, o := range liveroomdao.GetOnlinesByRoom(req.RoomId) {
		push.Data(o.UserId, cmd.LiveRoomGift, payload)
	}

	lockName := fmt.Sprintf("send_gift_%v", req.RoomId)
	gmlock.Lock(lockName)
	defer gmlock.Unlock(lockName)

	//防止并发,主播可以收到多个人的礼物
	liveRecord := liveroomdao.GetLiveRecordById(room.LiveRecordId)
	//添加本次直播收到的礼物总额
	liveRecord.AddTotalIncome(float64(totalCost))
	event.Pub(gameevent.RevenueEventEvent, eventData)

	//主播分成
	amount := float64(req.Count) * 0.4
	wallet.GoldAdd(req.RoomId, amount, currency.ReasonAnchorGiftRevenue)

	return &liveroomdto.SendGiftRes{
		Cost:    totalCost,
		Diamond: remaining,
	}, nil
}
