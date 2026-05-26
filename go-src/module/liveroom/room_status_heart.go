package liveroom

import (
	"context"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/os/gctx"
	"time"
	"xr-game-server/core/xrtimer"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
)

const (
	MaxFailNum   = 60 * 5
	StartFailNum = 60
)

var taskMap = gset.New(true)

func initHeart() {
	xrtimer.AddSingleton(gctx.New(), time.Second, heart)
}

// ReportLiveStartStatus 主播开播时上报开播状态
// 前端无需传参,后续在此补充开播记录等业务逻辑
func ReportLiveStartStatus(ctx context.Context, _ *liveroomdto.ReportLiveStartStatusReq) (*liveroomdto.ReportLiveStartStatusRes, error) {
	//先重置失败次数，防止被判断下播了
	room, _ := loadOwnRoom(ctx)
	if room.LiveRecordId == 0 {
		//没有开播
		return &liveroomdto.ReportLiveStartStatusRes{Success: true}, nil
	}
	room.ClearFailNum()
	liveRecord := liveroomdao.GetLiveRecordById(room.LiveRecordId)
	liveRecord.SetEndTime(nil)
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
	//如果达到1分钟不上报,直播计时停止
	if room.FailNum >= StartFailNum {
		if room.LiveRecordId > 0 {
			liveRecord := liveroomdao.GetLiveRecordById(room.LiveRecordId)
			if liveRecord.EndTime == nil {
				now := time.Now()
				liveRecord.SetEndTime(&now)
			}
		}
	}
	//如果达到最大值
	if room.FailNum >= MaxFailNum {
		//开始停止直播
		stopLive(userId)
		taskMap.Remove(userId)
		return
	}
	room.AddFailNum()
}
