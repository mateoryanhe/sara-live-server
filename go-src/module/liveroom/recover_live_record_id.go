package liveroom

import (
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/entity/live"
)

// 热重启后内存 live_record_id 可能仍为 0,而库已有绑定/未结束场次.
// 仅主播主动上报(心跳/下播)时补绑;观众心跳不补.查 live_rooms + 未结束 live_records;
// 心跳单房间最多查库 20 次,下播补绑不计入次数.
const maxLiveRecordRecoverDBAttempts = 20

var (
	liveRecordRecoverMu       sync.Mutex
	liveRecordRecoverAttempts = map[uint64]int{}
)

func resetLiveRecordRecoverAttempts(roomId uint64) {
	if roomId == 0 {
		return
	}
	liveRecordRecoverMu.Lock()
	delete(liveRecordRecoverAttempts, roomId)
	liveRecordRecoverMu.Unlock()
}

// tryRecoverLiveRecordBinding 内存 LiveRecordId==0 时从库补绑.
// consumeAttempt=true 时计入次数上限(心跳);false 时不计入(下播尽量圆 end_time).
func tryRecoverLiveRecordBinding(room *entity.LiveRoom, consumeAttempt bool) bool {
	if room == nil || room.ID == 0 || room.LiveRecordId > 0 {
		return room != nil && room.LiveRecordId > 0
	}
	if consumeAttempt {
		n := incrLiveRecordRecoverAttempts(room.ID)
		if n > maxLiveRecordRecoverDBAttempts {
			return false
		}
	}

	liveRecordId := resolveRecoverableLiveRecordId(room.ID)
	if liveRecordId == 0 {
		return false
	}
	applyRecoveredLiveRecordId(room, liveRecordId)
	xrlog.DetailLog.Infof(gctx.New(),
		"recover live_record_id roomId=%d liveRecordId=%d consumeAttempt=%v",
		room.ID, liveRecordId, consumeAttempt)
	return true
}

func incrLiveRecordRecoverAttempts(roomId uint64) int {
	liveRecordRecoverMu.Lock()
	defer liveRecordRecoverMu.Unlock()
	n := liveRecordRecoverAttempts[roomId] + 1
	liveRecordRecoverAttempts[roomId] = n
	return n
}

func resolveRecoverableLiveRecordId(roomId uint64) uint64 {
	if dbRoom := liveroomdao.GetRoomFromDB(roomId); dbRoom != nil && dbRoom.LiveRecordId > 0 {
		if isRecoverableOpenLiveRecord(dbRoom.LiveRecordId, roomId) {
			return dbRoom.LiveRecordId
		}
	}
	openId := liveroomdao.FindLatestOpenLiveRecordIdByAnchor(roomId)
	if openId == 0 {
		return 0
	}
	if !isRecoverableOpenLiveRecord(openId, roomId) {
		return 0
	}
	return openId
}

func isRecoverableOpenLiveRecord(liveRecordId, anchorId uint64) bool {
	row := liveroomdao.LoadLiveRecordFromDB(liveRecordId)
	if row == nil || row.ID == 0 {
		return false
	}
	if row.AnchorId != anchorId {
		return false
	}
	if row.EndTime != nil {
		return false
	}
	return true
}

func applyRecoveredLiveRecordId(room *entity.LiveRoom, liveRecordId uint64) {
	if room == nil || liveRecordId == 0 {
		return
	}
	room.SetLiveRecordId(liveRecordId)
	if room.HeartTime == nil {
		now := time.Now()
		room.SetHeartTime(&now)
	}
	resetLiveCameraEnabled(room.ID)
	liveroomdao.FlushRoomCache(room)
	addTask(room.ID)
	resetLiveRecordRecoverAttempts(room.ID)
}
