package liveroom

import (
	"context"
	"strconv"
	"strings"
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
	"xr-game-server/module/aliyunmoderation"
	"xr-game-server/module/livecfg"
	"xr-game-server/module/upload"
	"xr-game-server/module/wallet"
)

// SendPaidDanmaku 直播间付费弹幕
//  1. 按直播配置单价扣减钻石;不足时按缺口自动金币兑换钻石后扣款
//  2. 向房间内全体在线用户推送 cmd.LiveRoomPaidDanmaku
func SendPaidDanmaku(ctx context.Context, req *liveroomdto.SendPaidDanmakuReq) (*liveroomdto.SendPaidDanmakuRes, error) {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errercode.CreateCode(errercode.SysError)
	}

	price := livecfg.GetPaidDanmakuPrice()
	if price <= 0 {
		return nil, errercode.CreateCode(errercode.LiveRoomPaidDanmakuDisabled)
	}

	senderId := httpserver.GetAuthId(ctx)

	room := liveroomdao.GetRoomById(req.RoomId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	if room.LiveRecordId == 0 {
		return nil, errercode.CreateCode(errercode.LiveRoomNotLive)
	}

	if err := aliyunmoderation.RequireTextCompliant(aliyunmoderation.SceneChat, content); err != nil {
		return nil, err
	}

	if senderId != req.RoomId {
		onlineId := entity.BuildLiveRoomOnlineId(senderId, req.RoomId)
		if online := liveroomdao.GetOnlineById(onlineId, senderId, req.RoomId); online != nil && online.Muted {
			return nil, errercode.CreateCode(errercode.LiveRoomChatMuted)
		}
	}

	remaining, err := wallet.DiamondSubWithGoldExchange(senderId, price, currency.ReasonPaidDanmaku)
	if err != nil {
		return nil, err
	}

	eventData := entity.NewLiveRevenueLogRecord(
		room.ID, room.LiveRecordId, senderId, room.ID, 0, 1, price, price, uint8(liverevenue.PaidDanmaku),
	)

	sender := userinfodao.GetUserInfoByUserId(senderId)
	payload := &liveroomdto.PaidDanmakuPushItem{
		RoomId:    strconv.FormatUint(req.RoomId, 10),
		SenderId:  strconv.FormatUint(senderId, 10),
		Content:   content,
		UnitPrice: price,
		Cost:      price,
		SentAt:    time.Now().Unix(),
	}
	if sender != nil {
		payload.SenderName = sender.Nickname
		payload.SenderAvatar = upload.ResolveAvatarUrlForUser(sender.ID, sender.Avatar)
		payload.VipLevel = sender.VipLevel
	}

	for _, o := range getOnline(req.RoomId) {
		push.Data(o, cmd.LiveRoomPaidDanmaku, payload)
	}

	if liveRecord := liveroomdao.GetLiveRecordById(room.LiveRecordId); liveRecord != nil {
		liveRecord.AddPaidDanmakuEarn(price)
	}
	liveroomdao.GetLiveRoomIncomeUnsettled(room.ID).AddPaidDanmakuEarn(price)
	liveroomdao.GetLiveRoomIncomeTotal(room.ID).AddPaidDanmakuEarn(price)
	liveroomdao.MirrorGuildPaidDanmakuEarn(room.ID, price)
	liveroomdao.MirrorDailyAnchorPaidDanmakuEarn(room.ID, time.Now(), price)

	event.Pub(gameevent.RevenueEventEvent, eventData)

	NotifyLiveRecordTotalIncome(room)

	onlineId := entity.BuildLiveRoomOnlineId(senderId, req.RoomId)
	if online := liveroomdao.GetOnlineById(onlineId, senderId, req.RoomId); online != nil {
		online.AddTotalReward(price)
	}
	refreshRoomAudienceCaches(req.RoomId)

	return &liveroomdto.SendPaidDanmakuRes{
		Success: true,
		Cost:    price,
		Diamond: remaining,
	}, nil
}
