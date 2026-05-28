package liveroom

import (
	"context"
	"time"
	"xr-game-server/constants/cmd"
	"xr-game-server/constants/currency"
	"xr-game-server/core/actor"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/push"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/gift"
	"xr-game-server/module/upload"
	"xr-game-server/module/wallet"
)

const (
	Max  = 5
	Size = 100
)

var actorMap = make(map[uint64]*actor.Actor)

func initGiftActor() {
	for i := 0; i < Max; i++ {
		val := actor.NewActor(Size)
		val.Start()
		actorMap[uint64(i)] = val
	}
}

type GiftActorData struct {
	totalCost uint64
	RoomId    uint64
}

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

	//记录礼物流水日志
	entity.NewLiveGiftLogRecord(room.ID, room.LiveRecordId, senderId, room.ID, req.GiftId, req.Count, giftItem.Price, totalCost)

	//防止并发,主播可以收到多个人的礼物
	actorGift := actorMap[room.ID/Max]
	actorGift.Send(&GiftActorData{
		RoomId:    room.ID,
		totalCost: totalCost,
	}, func(val any) {
		liveRecord := liveroomdao.GetLiveRecordById(room.LiveRecordId)
		//添加本次直播收到的礼物总额
		liveRecord.AddTotalIncome(float64(totalCost))
	})

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

	return &liveroomdto.SendGiftRes{
		Cost:    totalCost,
		Diamond: remaining,
	}, nil
}
