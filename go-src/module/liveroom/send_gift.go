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
	userentity "xr-game-server/entity/user"
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

	roomGiftPayload := giftPushPayloadForRoomAudience(result.payload)
	for _, o := range getOnline(req.RoomId) {
		push.Data(o, cmd.LiveRoomGift, roomGiftPayload)
	}
	push.Data(req.RoomId, cmd.LiveRoomGift, roomGiftPayload)

	pushLuckyGiftWinByTier(result.senderId, req.RoomId, result.payload)
	recordSendGiftStats(result)
	return buildSendGiftRes(result), nil
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

	directPayload := buildDirectGiftPushItem(result.payload, anchorId)
	roomDirectPayload := directGiftPushPayloadForRoomAudience(directPayload)
	push.Data(result.senderId, cmd.LiveRoomDirectGift, directPayload)
	push.Data(anchorId, cmd.LiveRoomDirectGift, roomDirectPayload)

	pushLuckyGiftWinByTier(result.senderId, anchorId, result.payload)
	recordSendGiftStats(result)
	return buildSendGiftToAnchorRes(result), nil
}

type sendGiftResult struct {
	senderId  uint64
	room      *entity.LiveRoom
	giftId    uint64
	totalCost float64
	giftEarn  float64
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

	revenueType := uint8(liverevenue.Gift)
	giftEarn := totalCost
	var luckyDraw gift.LuckyGiftDrawResult

	if gift.IsLuckyGiftCategory(giftItem.CategoryType) {
		revenueType = uint8(liverevenue.GiftLucky)
		luckyDraw = gift.ResolveLuckyGiftPayout(giftId, totalCost)
		giftEarn = gift.CalcLuckyGiftAnchorSettlementEarn(totalCost)
		if luckyDraw.Payout > 0 {
			remaining, err = wallet.DiamondAdd(senderId, luckyDraw.Payout, currency.ReasonLuckyGiftPayout)
			if err != nil {
				return nil, err
			}
		}
	}

	eventData := entity.NewLiveRevenueLogRecord(room.ID, room.LiveRecordId, senderId, giftId, count, giftItem.Price, totalCost, revenueType)
	if revenueType == uint8(liverevenue.GiftLucky) {
		eventData.ApplyLuckyGiftMeta(luckyDraw.Rtp, luckyDraw.Multiplier, luckyDraw.Payout, giftEarn)
	}

	sender := userinfodao.GetUserInfoByUserId(senderId)
	payload := &liveroomdto.GiftPushItem{
		RoomId:              roomId,
		SenderId:            senderId,
		GiftId:              giftItem.ID,
		GiftName:            giftItem.Name,
		GiftIcon:            giftItem.Icon,
		GiftAnim:            giftItem.Animation,
		UnitPrice:           giftItem.Price,
		Count:               count,
		TotalCost:           totalCost,
		SentAt:              time.Now().Unix(),
		CategoryType:        giftItem.CategoryType,
		LuckyGiftRtp:        luckyDraw.Rtp,
		LuckyGiftMultiplier: luckyDraw.Multiplier,
		LuckyGiftPayout:     luckyDraw.Payout,
		LuckyGiftTierType:   luckyDraw.TierType,
	}
	applyGiftSenderProfile(payload, senderId, sender)

	return &sendGiftResult{
		senderId:  senderId,
		room:      room,
		giftId:    giftId,
		totalCost: totalCost,
		giftEarn:  giftEarn,
		remaining: remaining,
		payload:   payload,
		eventData: eventData,
	}, nil
}

func recordSendGiftStats(result *sendGiftResult) {
	if result == nil || result.room == nil {
		return
	}

	giftEarn := result.giftEarn
	if giftEarn <= 0 {
		giftEarn = result.totalCost
	}

	if liveRecord := liveroomdao.GetLiveRecordById(result.room.LiveRecordId); liveRecord != nil {
		liveRecord.AddGiftEarn(giftEarn)
		if result.room.LiveRecordId > 0 && liveroomdao.TryRecordLiveRecordGiftSender(result.room.LiveRecordId, result.senderId) {
			liveRecord.AddTotalGiftSender(1)
		}
		liveroomdao.PublishLiveRecord(liveRecord)
	}
	liveroomdao.GetLiveRoomIncomeUnsettled(result.room.ID).AddGiftEarn(giftEarn)
	liveroomdao.GetLiveRoomIncomeTotal(result.room.ID).AddGiftEarn(giftEarn)
	liveroomdao.MirrorGuildGiftEarn(result.room.ID, giftEarn)
	liveroomdao.MirrorDailyAnchorGiftEarn(result.room.ID, time.Now(), giftEarn)

	event.Pub(gameevent.RevenueEventEvent, result.eventData)
	liveroomdao.PublishRevenueLog(result.eventData)

	NotifyLiveRecordTotalIncome(result.room)

	onlineId := entity.BuildLiveRoomOnlineId(result.senderId, result.room.ID)
	if online := liveroomdao.GetOnlineById(onlineId, result.senderId, result.room.ID); online != nil {
		online.AddTotalReward(giftEarn)
		liveroomdao.PublishLiveRoomOnline(online)
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

func applyGiftSenderProfile(payload *liveroomdto.GiftPushItem, senderId uint64, sender *userentity.UserInfo) {
	if payload == nil || senderId == 0 {
		return
	}
	payload.SenderId = senderId
	brief := &liveroomdto.GiftSenderBrief{UserId: senderId}
	if sender != nil {
		brief.Nickname = sender.Nickname
		brief.Avatar = upload.ResolveAvatarUrlForUser(senderId, sender.Avatar)
		brief.VipLevel = sender.VipLevel
		payload.SenderName = brief.Nickname
		payload.SenderAvatar = brief.Avatar
		payload.VipLevel = brief.VipLevel
	}
	payload.Sender = brief
}

func buildSendGiftRes(result *sendGiftResult) *liveroomdto.SendGiftRes {
	if result == nil {
		return &liveroomdto.SendGiftRes{}
	}
	res := &liveroomdto.SendGiftRes{
		Cost:    result.totalCost,
		Diamond: result.remaining,
	}
	if result.payload != nil {
		res.Sender = result.payload.Sender
		res.LuckyGiftRtp = result.payload.LuckyGiftRtp
		res.LuckyGiftPayout = result.payload.LuckyGiftPayout
		res.LuckyGiftMultiplier = result.payload.LuckyGiftMultiplier
		res.LuckyGiftTierType = result.payload.LuckyGiftTierType
	}
	return res
}

func buildSendGiftToAnchorRes(result *sendGiftResult) *liveroomdto.SendGiftToAnchorRes {
	if result == nil {
		return &liveroomdto.SendGiftToAnchorRes{}
	}
	res := &liveroomdto.SendGiftToAnchorRes{
		Cost:    result.totalCost,
		Diamond: result.remaining,
	}
	if result.payload != nil {
		res.Sender = result.payload.Sender
		res.LuckyGiftRtp = result.payload.LuckyGiftRtp
		res.LuckyGiftPayout = result.payload.LuckyGiftPayout
		res.LuckyGiftMultiplier = result.payload.LuckyGiftMultiplier
		res.LuckyGiftTierType = result.payload.LuckyGiftTierType
	}
	return res
}

func buildLuckyGiftPayoutPushItem(payload *liveroomdto.GiftPushItem) *liveroomdto.LuckyGiftPayoutPushItem {
	if payload == nil {
		return nil
	}
	return &liveroomdto.LuckyGiftPayoutPushItem{
		RoomId:              payload.RoomId,
		Sender:              payload.Sender,
		GiftId:              payload.GiftId,
		GiftName:            payload.GiftName,
		GiftIcon:            payload.GiftIcon,
		Count:               payload.Count,
		TotalCost:           payload.TotalCost,
		LuckyGiftRtp:        payload.LuckyGiftRtp,
		LuckyGiftMultiplier: payload.LuckyGiftMultiplier,
		LuckyGiftPayout:     payload.LuckyGiftPayout,
		LuckyGiftTierType:   payload.LuckyGiftTierType,
		SentAt:              payload.SentAt,
	}
}

// pushLuckyGiftWinByTier 按 CMS 挡位 type 路由返奖展示推送：low 仅本人，medium 本房，high 全服
func pushLuckyGiftWinByTier(senderId, roomId uint64, payload *liveroomdto.GiftPushItem) {
	if senderId == 0 || payload == nil || payload.LuckyGiftPayout <= 0 {
		return
	}
	item := buildLuckyGiftPayoutPushItem(payload)
	if item == nil {
		return
	}
	switch payload.LuckyGiftTierType {
	case entity.LuckyGiftMultiplierTypeHigh:
		push.Broadcast(cmd.LiveRoomLuckyGiftPayout, item)
	case entity.LuckyGiftMultiplierTypeMedium:
		pushLuckyGiftPayoutToRoom(roomId, item)
	default:
		push.Data(senderId, cmd.LiveRoomLuckyGiftPayout, item)
	}
}

func pushLuckyGiftPayoutToRoom(roomId uint64, item *liveroomdto.LuckyGiftPayoutPushItem) {
	if roomId == 0 || item == nil {
		return
	}
	for _, o := range getOnline(roomId) {
		push.Data(o, cmd.LiveRoomLuckyGiftPayout, item)
	}
	push.Data(roomId, cmd.LiveRoomLuckyGiftPayout, item)
}

// giftPushPayloadForRoomAudience 低倍返奖仅本人可见：房间 cmd=8 不带返奖明细
func giftPushPayloadForRoomAudience(src *liveroomdto.GiftPushItem) *liveroomdto.GiftPushItem {
	if src == nil {
		return nil
	}
	out := *src
	if out.LuckyGiftPayout > 0 && out.LuckyGiftTierType == entity.LuckyGiftMultiplierTypeLow {
		out.LuckyGiftMultiplier = 0
		out.LuckyGiftPayout = 0
		out.LuckyGiftTierType = ""
	}
	return &out
}

func directGiftPushPayloadForRoomAudience(src *liveroomdto.DirectGiftPushItem) *liveroomdto.DirectGiftPushItem {
	if src == nil {
		return nil
	}
	out := *src
	if out.LuckyGiftPayout > 0 && out.LuckyGiftTierType == entity.LuckyGiftMultiplierTypeLow {
		out.LuckyGiftMultiplier = 0
		out.LuckyGiftPayout = 0
		out.LuckyGiftTierType = ""
	}
	return &out
}

func buildDirectGiftPushItem(gift *liveroomdto.GiftPushItem, anchorId uint64) *liveroomdto.DirectGiftPushItem {
	if gift == nil {
		return nil
	}
	return &liveroomdto.DirectGiftPushItem{
		RoomId:       gift.RoomId,
		AnchorId:     anchorId,
		SenderId:     gift.SenderId,
		SenderName:   gift.SenderName,
		SenderAvatar: gift.SenderAvatar,
		VipLevel:     gift.VipLevel,
		Sender:       gift.Sender,
		GiftId:       gift.GiftId,
		GiftName:     gift.GiftName,
		GiftIcon:     gift.GiftIcon,
		GiftAnim:     gift.GiftAnim,
		UnitPrice:           gift.UnitPrice,
		Count:               gift.Count,
		TotalCost:           gift.TotalCost,
		SentAt:              gift.SentAt,
		CategoryType:        gift.CategoryType,
		LuckyGiftRtp:          gift.LuckyGiftRtp,
		LuckyGiftMultiplier:   gift.LuckyGiftMultiplier,
		LuckyGiftPayout:       gift.LuckyGiftPayout,
		LuckyGiftTierType:     gift.LuckyGiftTierType,
	}
}
