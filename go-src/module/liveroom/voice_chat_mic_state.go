package liveroom

import (
	"sync"

	"xr-game-server/dao/liveroomdao"
	liveentity "xr-game-server/entity/live"
)

var voiceChatMicRoomLocks sync.Map // roomId -> *sync.Mutex

func voiceChatMicRoomLock(roomId uint64) *sync.Mutex {
	if roomId == 0 {
		return &sync.Mutex{}
	}
	v, _ := voiceChatMicRoomLocks.LoadOrStore(roomId, &sync.Mutex{})
	return v.(*sync.Mutex)
}

func voiceChatMicLiveRecordId(roomId uint64) uint64 {
	room := liveroomdao.GetRoomById(roomId)
	if room == nil {
		return 0
	}
	return room.LiveRecordId
}

func voiceChatMicOnStartLive(roomId, anchorId uint64) {
	liveRecordId := voiceChatMicLiveRecordId(roomId)
	if liveRecordId == 0 {
		return
	}
	liveroomdao.InitVoiceChatMicSession(roomId, liveRecordId, anchorId)
}

func voiceChatMicOnStopLive(roomId, liveRecordId uint64) {
	liveroomdao.ClearVoiceChatMicSession(roomId, liveRecordId)
	voiceChatMicRoomLocks.Delete(roomId)
}

func voiceChatMicOnUserLeave(roomId, userId uint64) {
	if roomId == 0 || userId == 0 {
		return
	}
	liveRecordId := voiceChatMicLiveRecordId(roomId)
	if liveRecordId == 0 {
		return
	}
	changed := false
	mu := voiceChatMicRoomLock(roomId)
	mu.Lock()
	for i := 0; i < liveentity.VoiceChatSeatCount; i++ {
		if i == liveentity.VoiceChatHostSeatIndex {
			continue
		}
		seat := liveroomdao.GetVoiceChatMicSeat(roomId, i)
		if seat == nil || seat.LiveRecordId != liveRecordId || seat.UserId != userId {
			continue
		}
		seat.SetUserId(0)
		seat.SetStatus(liveentity.VoiceChatMicSeatStatusIdle)
		liveroomdao.PublishVoiceChatMicSeat(seat)
		changed = true
		break
	}
	mu.Unlock()
	liveroomdao.CancelVoiceChatMicApply(roomId, userId)
	if changed {
		broadcastVoiceChatMicState(roomId)
	}
}

func voiceChatMicUserOnSeat(roomId, userId uint64) (seatIndex int, ok bool) {
	liveRecordId := voiceChatMicLiveRecordId(roomId)
	if liveRecordId == 0 || userId == 0 {
		return 0, false
	}
	for i := 0; i < liveentity.VoiceChatSeatCount; i++ {
		seat := liveroomdao.GetVoiceChatMicSeat(roomId, i)
		if seat != nil && seat.LiveRecordId == liveRecordId && seat.UserId == userId {
			return i, true
		}
	}
	return 0, false
}

func voiceChatMicHasApply(roomId, userId uint64) bool {
	liveRecordId := voiceChatMicLiveRecordId(roomId)
	return liveroomdao.HasActiveVoiceChatMicApply(roomId, userId, liveRecordId)
}

func voiceChatMicSnapshot(roomId uint64) ([liveentity.VoiceChatSeatCount]uint64, []*liveentity.VoiceChatMicApply) {
	var seats [liveentity.VoiceChatSeatCount]uint64
	liveRecordId := voiceChatMicLiveRecordId(roomId)
	if liveRecordId == 0 {
		return seats, nil
	}
	for i := 0; i < liveentity.VoiceChatSeatCount; i++ {
		seat := liveroomdao.GetVoiceChatMicSeat(roomId, i)
		if seat != nil && seat.LiveRecordId == liveRecordId {
			seats[i] = seat.UserId
		}
	}
	return seats, liveroomdao.GetActiveVoiceChatMicApplies(roomId, liveRecordId)
}

func voiceChatMicFirstEmptyAudienceSeat(roomId uint64, liveRecordId uint64) int {
	for i := liveentity.VoiceChatHostSeatIndex + 1; i < liveentity.VoiceChatSeatCount; i++ {
		seat := liveroomdao.GetVoiceChatMicSeat(roomId, i)
		if seat != nil && seat.LiveRecordId == liveRecordId && seat.UserId == 0 &&
			liveentity.VoiceChatMicSeatStatusAllowsAudienceTake(seat.Status) {
			return i
		}
	}
	return -1
}

func voiceChatMicIsValidAudienceSeat(seatIndex int) bool {
	return seatIndex > liveentity.VoiceChatHostSeatIndex && seatIndex < liveentity.VoiceChatSeatCount
}

func voiceChatMicIsValidSeatIndex(seatIndex int) bool {
	return seatIndex >= liveentity.VoiceChatHostSeatIndex && seatIndex < liveentity.VoiceChatSeatCount
}

func voiceChatMicAssignAudienceSeat(roomId, userId uint64, preferredSeat int) (seatIndex int, ok bool) {
	liveRecordId := voiceChatMicLiveRecordId(roomId)
	if liveRecordId == 0 || userId == 0 {
		return 0, false
	}
	mu := voiceChatMicRoomLock(roomId)
	mu.Lock()
	defer mu.Unlock()

	if idx, on := voiceChatMicUserOnSeat(roomId, userId); on {
		return idx, true
	}

	seatIndex = preferredSeat
	if seatIndex == 0 {
		seatIndex = voiceChatMicFirstEmptyAudienceSeat(roomId, liveRecordId)
	} else if !voiceChatMicIsValidAudienceSeat(seatIndex) {
		return 0, false
	}
	if seatIndex < 0 {
		return 0, false
	}
	seat := liveroomdao.GetVoiceChatMicSeat(roomId, seatIndex)
	if seat == nil || seat.LiveRecordId != liveRecordId || seat.UserId != 0 || !liveentity.VoiceChatMicSeatStatusAllowsAudienceTake(seat.Status) {
		return 0, false
	}
	seat.SetUserId(userId)
	seat.SetStatus(liveentity.VoiceChatMicSeatStatusOnMic)
	liveroomdao.PublishVoiceChatMicSeat(seat)
	liveroomdao.CancelVoiceChatMicApply(roomId, userId)
	return seatIndex, true
}

func voiceChatMicSwitchAudienceSeat(roomId, userId uint64, targetSeat int) (fromSeat, toSeat int, ok bool) {
	liveRecordId := voiceChatMicLiveRecordId(roomId)
	if liveRecordId == 0 || userId == 0 {
		return 0, 0, false
	}
	if targetSeat != 0 && !voiceChatMicIsValidAudienceSeat(targetSeat) {
		return 0, 0, false
	}
	mu := voiceChatMicRoomLock(roomId)
	mu.Lock()
	defer mu.Unlock()

	fromSeat = -1
	for i := 0; i < liveentity.VoiceChatSeatCount; i++ {
		if i == liveentity.VoiceChatHostSeatIndex {
			continue
		}
		seat := liveroomdao.GetVoiceChatMicSeat(roomId, i)
		if seat != nil && seat.LiveRecordId == liveRecordId && seat.UserId == userId {
			fromSeat = i
			break
		}
	}
	if fromSeat < 0 {
		return 0, 0, false
	}
	if targetSeat == 0 {
		targetSeat = fromSeat
	}
	if targetSeat == fromSeat {
		return fromSeat, fromSeat, true
	}
	from := liveroomdao.GetVoiceChatMicSeat(roomId, fromSeat)
	to := liveroomdao.GetVoiceChatMicSeat(roomId, targetSeat)
	if from == nil || to == nil || from.LiveRecordId != liveRecordId || to.LiveRecordId != liveRecordId {
		return 0, 0, false
	}
	if to.UserId != 0 || !liveentity.VoiceChatMicSeatStatusAllowsAudienceTake(to.Status) {
		return 0, 0, false
	}
	newStatus := liveentity.VoiceChatMicSeatStatusOnMic
	if from.Status == liveentity.VoiceChatMicSeatStatusMicMuted {
		newStatus = liveentity.VoiceChatMicSeatStatusMicMuted
	}
	from.SetUserId(0)
	from.SetStatus(liveentity.VoiceChatMicSeatStatusIdle)
	liveroomdao.PublishVoiceChatMicSeat(from)
	to.SetUserId(userId)
	to.SetStatus(newStatus)
	liveroomdao.PublishVoiceChatMicSeat(to)
	return fromSeat, targetSeat, true
}

func voiceChatMicClearUserFromSeat(roomId, liveRecordId, userId uint64) bool {
	if roomId == 0 || userId == 0 {
		return false
	}
	mu := voiceChatMicRoomLock(roomId)
	mu.Lock()
	defer mu.Unlock()
	for i := 0; i < liveentity.VoiceChatSeatCount; i++ {
		if i == liveentity.VoiceChatHostSeatIndex {
			continue
		}
		seat := liveroomdao.GetVoiceChatMicSeat(roomId, i)
		if seat == nil || seat.LiveRecordId != liveRecordId || seat.UserId != userId {
			continue
		}
		seat.SetUserId(0)
		seat.SetStatus(liveentity.VoiceChatMicSeatStatusIdle)
		liveroomdao.PublishVoiceChatMicSeat(seat)
		return true
	}
	return false
}
