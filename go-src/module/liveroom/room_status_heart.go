package liveroom

import (
	"context"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/os/gctx"
	"time"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
)

const (
	TimeOut = 5 * time.Minute
)

var taskMap = gset.New(true)

func initHeart() {
	ids := liveroomdao.ListLivingRoomIds()
	for _, id := range ids {
		taskMap.Add(id)
	}
	xrtimer.AddSingleton(gctx.New(), time.Second, heart)
}

// ReportLiveStartStatus 主播开播时上报开播状态
// 前端无需传参,后续在此补充开播记录等业务逻辑,每秒上报一次,形成开播时间
func ReportLiveStartStatus(ctx context.Context, _ *liveroomdto.ReportLiveStartStatusReq) (*liveroomdto.ReportLiveStartStatusRes, error) {
	//先重置失败次数，防止被判断下播了
	room, _ := loadOwnRoom(ctx)
	if room.LiveRecordId == 0 {
		//没有开播
		return &liveroomdto.ReportLiveStartStatusRes{Success: true}, nil
	}
	//刷新房间状态,防止离线
	now := time.Now()
	room.SetHeartTime(&now)
	//记录本次直播
	liveRecord := liveroomdao.GetLiveRecordById(room.LiveRecordId)
	liveRecord.AddTotalLiveDuration(1)
	//累计全部直播
	stat := userinfodao.GetUserCumulativeStatByUserId(room.ID)
	stat.AddTotalLiveDuration(1)
	return &liveroomdto.ReportLiveStartStatusRes{Success: true}, nil
}

func addTask(userId uint64) {
	taskMap.Add(userId)
}

func heart(ctx context.Context) {
	if taskMap.Size() == 0 {
		return
	}
	//每秒检查一下直播间
	temp := make([]any, 0)
	temp = append(temp, taskMap.Slice())
	for _, data := range temp {
		userId := data.(uint64)
		chkOne(userId)
	}
}

func chkOne(userId uint64) {
	room := liveroomdao.GetRoomById(userId)
	now := time.Now()
	if room.HeartTime == nil {
		room.SetHeartTime(&now)
	}
	diff := now.Sub(*room.HeartTime)
	//如果相差超过5分钟,判定离线
	if TimeOut > diff {
		return
	}
	stopLive(userId)
}
