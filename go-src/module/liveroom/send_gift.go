package liveroom

import (
	"context"
	"math"
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
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/gameevent"
	"xr-game-server/module/gift"
	"xr-game-server/module/upload"
	"xr-game-server/module/wallet"
)

// SendGift 直播间送礼
//  1. 校验房间存在、礼物存在(命中礼物缓存,即默认已上架)、数量合法
//  2. 计算总消耗 = 礼物单价 * 数量,优先扣钻石;不足时按缺口自动金币兑换钻石后扣款
//  3. 扣款成功后,向房间内全体在线用户(含送礼人自身)推送 cmd.LiveRoomGift
func SendGift(ctx context.Context, req *liveroomdto.SendGiftReq) (*liveroomdto.SendGiftRes, error) {
	result, err := prepareSendGift(ctx, req.RoomId, req.GiftId, req.Count)
	if err != nil {
		return nil, err
	}

	for _, o := range getOnline(req.RoomId) {
		push.Data(o, cmd.LiveRoomGift, result.payload)
	}
	push.Data(req.RoomId, cmd.LiveRoomGift, result.payload)

	recordSendGiftStats(result)
	return &liveroomdto.SendGiftRes{
		Cost:    result.totalCost,
		Diamond: result.remaining,
	}, nil
}

// SendGiftToAnchor 给指定主播送礼(仅推送给发送方与主播)
func SendGiftToAnchor(ctx context.Context, req *liveroomdto.SendGiftToAnchorReq) (*liveroomdto.SendGiftToAnchorRes, error) {
	anchorId := req.AnchorId
	if anchorId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	result, err := prepareSendGift(ctx, anchorId, req.GiftId, req.Count)
	if err != nil {
		return nil, err
	}
	if result.senderId == anchorId {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	payload := buildPrivateGiftPushItem(result.payload, anchorId)
	push.Data(result.senderId, cmd.LiveRoomPrivateGift, payload)
	push.Data(anchorId, cmd.LiveRoomPrivateGift, payload)

	recordSendGiftStats(result)
	return &liveroomdto.SendGiftToAnchorRes{
		Cost:    result.totalCost,
		Diamond: result.remaining,
	}, nil
}

type sendGiftResult struct {
	senderId  uint64
	room      *entity.LiveRoom
	giftId    uint64
	totalCost float64
	remaining float64
	payload   *liveroomdto.GiftPushItem
	eventData *entity.LiveRevenueLog
}

func prepareSendGift(ctx context.Context, roomId, giftId uint64, count int) (*sendGiftResult, error) {
	if count <= 0 {
		return nil, errercode.CreateCode(errercode.GiftCountInvalid)
	}

	senderId := httpserver.GetAuthId(ctx)
	room := liveroomdao.GetRoomById(roomId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	giftItem := gift.GetGiftFromCacheById(giftId)
	if giftItem == nil {
		return nil, errercode.CreateCode(errercode.GiftOffShelf)
	}

	totalCost, err := calcSendGiftTotalCost(giftItem.Price, count)
	if err != nil {
		return nil, err
	}
	var remaining float64
	if totalCost > 0 {
		remaining, err = wallet.DiamondSubWithGoldExchange(senderId, totalCost, currency.ReasonGiftSend)
		if err != nil {
			return nil, err
		}
	} else {
		remaining = userinfodao.GetUserInfoByUserId(senderId).Diamond
	}

	eventData := entity.NewLiveRevenueLogRecord(room.ID, room.LiveRecordId, senderId, room.ID, giftId, count, giftItem.Price, totalCost, uint8(liverevenue.Gift))

	sender := userinfodao.GetUserInfoByUserId(senderId)
	payload := &liveroomdto.GiftPushItem{
		RoomId:    roomId,
		SenderId:  senderId,
		GiftId:    giftItem.ID,
		GiftName:  giftItem.Name,
		GiftIcon:  giftItem.Icon,
		GiftAnim:  giftItem.Animation,
		UnitPrice: giftItem.Price,
		Count:     count,
		TotalCost: totalCost,
		SentAt:    time.Now().Unix(),
	}
	if sender != nil {
		payload.SenderName = sender.Nickname
		payload.SenderAvatar = upload.ResolveAvatarUrlForUser(sender.ID, sender.Avatar)
		payload.VipLevel = sender.VipLevel
	}

	return &sendGiftResult{
		senderId:  senderId,
		room:      room,
		giftId:    giftId,
		totalCost: totalCost,
		remaining: remaining,
		payload:   payload,
		eventData: eventData,
	}, nil
}

func recordSendGiftStats(result *sendGiftResult) {
	if result == nil || result.room == nil {
		return
	}

	if liveRecord := liveroomdao.GetLiveRecordById(result.room.LiveRecordId); liveRecord != nil {
		liveRecord.AddGiftEarn(result.totalCost)
		if result.room.LiveRecordId > 0 && liveroomdao.TryRecordLiveRecordGiftSender(result.room.LiveRecordId, result.senderId) {
			liveRecord.AddTotalGiftSender(1)
		}
	}
	liveroomdao.GetLiveRoomIncomeUnsettled(result.room.ID).AddGiftEarn(result.totalCost)
	liveroomdao.GetLiveRoomIncomeTotal(result.room.ID).AddGiftEarn(result.totalCost)
	liveroomdao.MirrorGuildGiftEarn(result.room.ID, result.totalCost)

	event.Pub(gameevent.RevenueEventEvent, result.eventData)

	onlineId := entity.BuildLiveRoomOnlineId(result.senderId, result.room.ID)
	if online := liveroomdao.GetOnlineById(onlineId, result.senderId, result.room.ID); online != nil {
		online.AddTotalReward(result.totalCost)
	}
	refreshRoomAudienceCaches(result.room.ID)
}

func calcSendGiftTotalCost(unitPrice float64, count int) (float64, error) {
	if count <= 0 {
		return 0, errercode.CreateCode(errercode.GiftCountInvalid)
	}
	if unitPrice <= 0 {
		return 0, nil
	}
	if float64(count) > math.MaxFloat64/unitPrice {
		return 0, errercode.CreateCode(errercode.GiftCountInvalid)
	}
	total := unitPrice * float64(count)
	if total <= 0 || math.IsInf(total, 0) || math.IsNaN(total) {
		return 0, errercode.CreateCode(errercode.GiftCountInvalid)
	}
	return total, nil
}

func buildPrivateGiftPushItem(gift *liveroomdto.GiftPushItem, anchorId uint64) *liveroomdto.PrivateGiftPushItem {
	if gift == nil {
		return nil
	}
	return &liveroomdto.PrivateGiftPushItem{
		RoomId:       gift.RoomId,
		AnchorId:     anchorId,
		SenderId:     gift.SenderId,
		SenderName:   gift.SenderName,
		SenderAvatar: gift.SenderAvatar,
		VipLevel:     gift.VipLevel,
		GiftId:       gift.GiftId,
		GiftName:     gift.GiftName,
		GiftIcon:     gift.GiftIcon,
		GiftAnim:     gift.GiftAnim,
		UnitPrice:    gift.UnitPrice,
		Count:        gift.Count,
		TotalCost:    gift.TotalCost,
		SentAt:       gift.SentAt,
	}
}
