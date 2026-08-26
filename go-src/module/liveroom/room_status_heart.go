package liveroom

import (
	"context"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/os/gctx"
	"time"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
)

const (
	TimeOut              = 30 * time.Minute
	PrivatePeriod uint64 = 30
	CommonPeriod  uint64 = 30
)

var taskMap = gset.NewTSet[uint64](true)

func initHeart() {
	ids := liveroomdao.ListLivingRoomIds()
	for _, id := range ids {
		taskMap.Add(id)
	}
	allOnline := liveroomdao.GetAllOnline()
	for _, data := range allOnline {
		initRoomOnline(data.RoomId)
		addToOnline(data.UserId, data.RoomId)
	}
	for _, data := range ids {
		refreshRoomAudienceCaches(data)
	}
	xrtimer.AddOnce(gctx.New(), time.Minute, func(ctx context.Context) {
		xrtimer.AddSingleton(gctx.New(), time.Second, heart)
	})

}

// ReportLiveStartStatus 主播开播时上报开播状态
// 前端无需传参,后续在此补充开播记录等业务逻辑,每30秒上报一次,形成开播时间
func ReportLiveStartStatus(ctx context.Context, req *liveroomdto.ReportLiveStartStatusReq) (*liveroomdto.ReportLiveStartStatusRes, error) {
	userId := httpserver.GetAuthId(ctx)
	room := liveroomdao.GetRoomById(req.RoomId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	if room.LiveRecordId == 0 {
		return &liveroomdto.ReportLiveStartStatusRes{Success: true}, nil
	}

	var billingDeducted float64

	if userId == room.ID {
		flushAnchorId(room)
	} else {
		//检查是否加入房间
		if !isUserInOnlineMap(userId, room.ID) {
			return &liveroomdto.ReportLiveStartStatusRes{Success: true}, nil
		}
		deducted, err := chargePrivateRoomBillingIfNeeded(userId, room)
		if err != nil {
			exitRoom(userId, room.ID)
			return nil, err
		}
		billingDeducted = deducted
		flushAudience(userId, room)
	}

	return &liveroomdto.ReportLiveStartStatusRes{
		Success:         true,
		BillingDeducted: billingDeducted,
	}, nil
}

func flushAnchorId(room *entity.LiveRoom) {
	//刷新房间状态,防止离线
	now := time.Now()
	room.SetHeartTime(&now)
	liveroomdao.FlushRoomCache(room)
	period := CommonPeriod
	if cfg := liveroomdao.GetLiveRoomCfg(room.ID); cfg != nil && cfg.Category == entity.LiveRoomCategoryPrivate {
		period = PrivatePeriod
	}
	sec := float64(period)
	//记录本次直播
	liveRecord := liveroomdao.GetLiveRecordById(room.LiveRecordId)
	liveRecord.AddTotalLiveDuration(sec)
	liveroomdao.PublishLiveRecord(liveRecord)
	//直播间未结算+生涯累计时长
	if unsettled := liveroomdao.GetLiveRoomIncomeUnsettled(room.ID); unsettled != nil {
		unsettled.AddTotalLiveDuration(sec)
	}
	if total := liveroomdao.GetLiveRoomIncomeTotal(room.ID); total != nil {
		total.AddTotalLiveDuration(sec)
	}
	liveroomdao.MirrorGuildLiveDuration(room.ID, sec)
	liveroomdao.MirrorDailyAnchorLiveDuration(room.ID, now, sec)
}

func flushAudience(userId uint64, room *entity.LiveRoom) {
	timeNow := time.Now()
	onlineId := entity.BuildLiveRoomOnlineId(userId, room.ID)
	onlineData := liveroomdao.GetOnlineById(onlineId, userId, room.ID)
	if onlineData == nil {
		return
	}
	onlineData.SetHeartTime(&timeNow)
	liveroomdao.PublishLiveRoomOnline(onlineData)
}

func addTask(userId uint64) {
	taskMap.Add(userId)
}

func heart(ctx context.Context) {
	if taskMap.Size() == 0 {
		return
	}
	//每秒检查一下直播间
	for _, roomId := range taskMap.Slice() {
		chkAnchor(roomId)
		chkAudience(roomId)
	}
}

func chkAnchor(roomId uint64) {
	room := liveroomdao.GetRoomById(roomId)
	now := time.Now()
	if room.HeartTime == nil {
		room.SetHeartTime(&now)
	}
	diff := now.Sub(*room.HeartTime)
	//如果相差超过5分钟,判定离线
	if TimeOut > diff {
		return
	}
	stopLive(roomId)
}

func chkAudience(roomId uint64) {
	allOnline := getOnline(roomId)
	now := time.Now()
	for _, data := range allOnline {
		onlineId := entity.BuildLiveRoomOnlineId(data, roomId)
		onlineData := liveroomdao.GetOnlineById(onlineId, data, roomId)
		if onlineData.HeartTime == nil {
			onlineData.SetHeartTime(&now)
			liveroomdao.PublishLiveRoomOnline(onlineData)
		}
		diff := now.Sub(*onlineData.HeartTime)
		//如果相差超过5分钟,判定离线
		if TimeOut > diff {
			return
		}
		exitRoom(data, roomId)
	}
}
