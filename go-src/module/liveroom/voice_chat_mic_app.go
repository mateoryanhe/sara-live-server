package liveroom

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/constants/cmd"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/push"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/liveroomdto"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

func voiceChatMicSwitchAudienceSeatErr(roomId, liveRecordId uint64, targetSeat int) errercode.XRCode {
	if !voiceChatMicIsValidAudienceSeat(targetSeat) {
		return errercode.VoiceChatMicSeatInvalid
	}
	seat := liveroomdao.GetVoiceChatMicSeat(roomId, targetSeat)
	if seat != nil && seat.LiveRecordId == liveRecordId && seat.UserId == 0 && liveentity.IsVoiceChatMicSeatLocked(seat) {
		return errercode.VoiceChatMicSeatLocked
	}
	return errercode.VoiceChatMicSeatTaken
}

func dedupeAudienceSeatIndexes(indexes []int) []int {
	if len(indexes) == 0 {
		return nil
	}
	seen := make(map[int]struct{}, len(indexes))
	out := make([]int, 0, len(indexes))
	for _, idx := range indexes {
		if !voiceChatMicIsValidAudienceSeat(idx) {
			continue
		}
		if _, ok := seen[idx]; ok {
			continue
		}
		seen[idx] = struct{}{}
		out = append(out, idx)
	}
	return out
}

func voiceChatMicAssignAudienceSeatErr(roomId, liveRecordId uint64, preferredSeat int) errercode.XRCode {
	if preferredSeat != 0 && voiceChatMicIsValidAudienceSeat(preferredSeat) {
		seat := liveroomdao.GetVoiceChatMicSeat(roomId, preferredSeat)
		if seat != nil && seat.LiveRecordId == liveRecordId && seat.UserId == 0 && liveentity.IsVoiceChatMicSeatLocked(seat) {
			return errercode.VoiceChatMicSeatLocked
		}
	}
	if preferredSeat == 0 {
		return errercode.VoiceChatMicSeatFull
	}
	return errercode.VoiceChatMicSeatTaken
}

func isVoiceChatRoomId(roomId uint64) bool {
	if roomId == 0 {
		return false
	}
	cfg := liveroomdao.GetLiveRoomCfgFromCache(roomId)
	if cfg == nil {
		cfg = liveroomdao.GetLiveRoomCfg(roomId)
	}
	if cfg == nil {
		return false
	}
	return normalizeLiveRoomCategory(cfg.Category) == liveentity.LiveRoomCategoryVoiceChat
}

func loadVoiceChatLiveRoom(roomId uint64) (*liveentity.LiveRoom, error) {
	room := liveroomdao.GetRoomById(roomId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	if IsRoomBanned(room) {
		return nil, errercode.CreateCode(errercode.LiveRoomBanned)
	}
	if !isVoiceChatRoomId(roomId) {
		return nil, errercode.CreateCode(errercode.LiveRoomNotVoiceChat)
	}
	if room.LiveRecordId == 0 {
		return nil, errercode.CreateCode(errercode.LiveRoomNotLive)
	}
	return room, nil
}

func requireAudienceOnline(roomId, userId uint64) error {
	if userId == roomId {
		return nil
	}
	onlineId := liveentity.BuildLiveRoomOnlineId(userId, roomId)
	online := liveroomdao.GetOnlineById(onlineId, userId, roomId)
	if online == nil || online.Status != liveentity.LiveRoomOnlineStatusOnline {
		return errercode.CreateCode(errercode.LiveRoomAudienceNotOnline)
	}
	return nil
}

func toVoiceChatMicUserBrief(userId uint64) *liveroomdto.VoiceChatMicUserBrief {
	if userId == 0 {
		return nil
	}
	u := userinfodao.GetUserInfoByUserId(userId)
	if u == nil {
		return &liveroomdto.VoiceChatMicUserBrief{
			UserId: strconv.FormatUint(userId, 10),
		}
	}
	return &liveroomdto.VoiceChatMicUserBrief{
		UserId:   strconv.FormatUint(userId, 10),
		Nickname: u.Nickname,
		Avatar:   upload.ResolveAvatarUrlForUser(userId, u.Avatar),
		VipLevel: u.VipLevel,
		Gender:   u.Gender,
	}
}

func buildVoiceChatMicStatePayload(roomId uint64) *liveroomdto.VoiceChatMicStatePayload {
	seats, _ := voiceChatMicSnapshot(roomId)
	items := make([]*liveroomdto.VoiceChatMicSeatItem, 0, liveentity.VoiceChatSeatCount)
	liveRecordId := voiceChatMicLiveRecordId(roomId)
	for i := 0; i < liveentity.VoiceChatSeatCount; i++ {
		userId := seats[i]
		status := liveentity.VoiceChatMicSeatStatusIdle
		seat := liveroomdao.GetVoiceChatMicSeat(roomId, i)
		if seat != nil && seat.LiveRecordId == liveRecordId {
			status = seat.Status
		}
		item := &liveroomdto.VoiceChatMicSeatItem{
			SeatIndex: i,
			Status:    status,
			Occupied:  userId > 0,
		}
		if userId > 0 {
			item.User = toVoiceChatMicUserBrief(userId)
		}
		items = append(items, item)
	}
	return &liveroomdto.VoiceChatMicStatePayload{
		RoomId:           strconv.FormatUint(roomId, 10),
		VoiceChatMicMode: voiceChatMicModeForRoom(roomId),
		UpdatedAt:        time.Now().UnixMilli(),
		Seats:            items,
	}
}

func broadcastVoiceChatMicState(roomId uint64) {
	if roomId == 0 {
		return
	}
	payload := buildVoiceChatMicStatePayload(roomId)
	PushToRoomAudience(roomId, cmd.LiveRoomVoiceChatMicState, payload)
	push.Data(roomId, cmd.LiveRoomVoiceChatMicState, payload)
}

// GetVoiceChatMicState 查询语聊房麦位状态
func GetVoiceChatMicState(ctx context.Context, req *liveroomdto.GetVoiceChatMicStateReq) (*liveroomdto.GetVoiceChatMicStateRes, error) {
	if _, err := loadVoiceChatLiveRoom(req.RoomId); err != nil {
		return nil, err
	}
	return buildVoiceChatMicStatePayload(req.RoomId), nil
}

// ApplyVoiceChatMic 观众申请上麦
func ApplyVoiceChatMic(ctx context.Context, req *liveroomdto.ApplyVoiceChatMicReq) (*liveroomdto.ApplyVoiceChatMicRes, error) {
	userId := httpserver.GetAuthId(ctx)
	room, err := loadVoiceChatLiveRoom(req.RoomId)
	if err != nil {
		return nil, err
	}
	if userId == room.ID {
		return nil, errercode.CreateCode(errercode.VoiceChatMicHostNoApply)
	}
	if err := requireAudienceOnline(room.ID, userId); err != nil {
		return nil, err
	}
	if _, onSeat := voiceChatMicUserOnSeat(room.ID, userId); onSeat {
		return nil, errercode.CreateCode(errercode.VoiceChatMicAlreadyOn)
	}
	micMode := voiceChatMicModeForRoom(room.ID)
	switch micMode {
	case liveentity.VoiceChatMicModeHostOnly:
		return nil, errercode.CreateCode(errercode.VoiceChatMicModeHostOnly)
	case liveentity.VoiceChatMicModeFree:
		if voiceChatMicHasApply(room.ID, userId) {
			return nil, errercode.CreateCode(errercode.VoiceChatMicApplyExist)
		}
		preferredSeat := req.SeatIndex
		if preferredSeat != 0 && !voiceChatMicIsValidAudienceSeat(preferredSeat) {
			return nil, errercode.CreateCode(errercode.VoiceChatMicSeatInvalid)
		}
		seatIndex, ok := voiceChatMicAssignAudienceSeat(room.ID, userId, preferredSeat)
		if !ok {
			return nil, errercode.CreateCode(voiceChatMicAssignAudienceSeatErr(room.ID, room.LiveRecordId, preferredSeat))
		}
		broadcastVoiceChatMicState(room.ID)
		return &liveroomdto.ApplyVoiceChatMicRes{Success: true, SeatIndex: seatIndex}, nil
	}

	if voiceChatMicHasApply(room.ID, userId) {
		return nil, errercode.CreateCode(errercode.VoiceChatMicApplyExist)
	}
	if liveroomdao.ActiveVoiceChatMicApplyCount(room.ID, room.LiveRecordId) >= liveentity.VoiceChatMicApplyListMax {
		return nil, errercode.CreateCode(errercode.VoiceChatMicApplyListFull)
	}
	if voiceChatMicFirstEmptyAudienceSeat(room.ID, room.LiveRecordId) < 0 {
		return nil, errercode.CreateCode(errercode.VoiceChatMicSeatFull)
	}
	liveroomdao.UpsertActiveVoiceChatMicApply(room.ID, userId, room.LiveRecordId, time.Now())
	return &liveroomdto.ApplyVoiceChatMicRes{Success: true}, nil
}

// CancelVoiceChatMicApply 取消上麦申请
func CancelVoiceChatMicApply(ctx context.Context, req *liveroomdto.CancelVoiceChatMicApplyReq) (*liveroomdto.CancelVoiceChatMicApplyRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if _, err := loadVoiceChatLiveRoom(req.RoomId); err != nil {
		return nil, err
	}
	liveroomdao.CancelVoiceChatMicApply(req.RoomId, userId)
	return &liveroomdto.CancelVoiceChatMicApplyRes{Success: true}, nil
}

// VoiceChatMicApplyList 主播查询上麦申请列表
func VoiceChatMicApplyList(ctx context.Context, req *liveroomdto.VoiceChatMicApplyListReq) (*liveroomdto.VoiceChatMicApplyListRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	room, err := loadVoiceChatLiveRoom(req.RoomId)
	if err != nil {
		return nil, err
	}
	if anchorId != room.ID {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	if voiceChatMicModeForRoom(room.ID) != liveentity.VoiceChatMicModeApply {
		return &liveroomdto.VoiceChatMicApplyListRes{
			RoomId: strconv.FormatUint(room.ID, 10),
			Total:  0,
			List:   make([]*liveroomdto.VoiceChatMicApplyItem, 0),
		}, nil
	}

	applications := liveroomdao.GetActiveVoiceChatMicApplies(room.ID, room.LiveRecordId)
	list := make([]*liveroomdto.VoiceChatMicApplyItem, 0, len(applications))
	for _, item := range applications {
		if item == nil {
			continue
		}
		appliedAt := int64(0)
		if item.AppliedAt != nil {
			appliedAt = item.AppliedAt.Unix()
		}
		list = append(list, &liveroomdto.VoiceChatMicApplyItem{
			User:      toVoiceChatMicUserBrief(item.UserId),
			AppliedAt: appliedAt,
		})
	}
	if list == nil {
		list = make([]*liveroomdto.VoiceChatMicApplyItem, 0)
	}
	return &liveroomdto.VoiceChatMicApplyListRes{
		RoomId: strconv.FormatUint(room.ID, 10),
		Total:  len(list),
		List:   list,
	}, nil
}

// ApproveVoiceChatMicApply 主播同意上麦
func ApproveVoiceChatMicApply(ctx context.Context, req *liveroomdto.ApproveVoiceChatMicApplyReq) (*liveroomdto.ApproveVoiceChatMicApplyRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	room, err := loadVoiceChatLiveRoom(req.RoomId)
	if err != nil {
		return nil, err
	}
	if anchorId != room.ID {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	if voiceChatMicModeForRoom(room.ID) != liveentity.VoiceChatMicModeApply {
		return nil, errercode.CreateCode(errercode.VoiceChatMicModeNotApply)
	}
	if req.UserId == room.ID {
		return nil, errercode.CreateCode(errercode.VoiceChatMicHostSeatLocked)
	}
	if err := requireAudienceOnline(room.ID, req.UserId); err != nil {
		return nil, err
	}

	if !voiceChatMicHasApply(room.ID, req.UserId) {
		return nil, errercode.CreateCode(errercode.VoiceChatMicApplyNotFound)
	}
	if _, onSeat := voiceChatMicUserOnSeat(room.ID, req.UserId); onSeat {
		return nil, errercode.CreateCode(errercode.VoiceChatMicAlreadyOn)
	}
	preferredSeat := req.SeatIndex
	if preferredSeat != 0 && !voiceChatMicIsValidAudienceSeat(preferredSeat) {
		return nil, errercode.CreateCode(errercode.VoiceChatMicSeatInvalid)
	}
	seatIndex, ok := voiceChatMicAssignAudienceSeat(room.ID, req.UserId, preferredSeat)
	if !ok {
		return nil, errercode.CreateCode(voiceChatMicAssignAudienceSeatErr(room.ID, room.LiveRecordId, preferredSeat))
	}
	broadcastVoiceChatMicState(room.ID)
	return &liveroomdto.ApproveVoiceChatMicApplyRes{Success: true, SeatIndex: seatIndex}, nil
}

func setVoiceChatMicSeatLockedByAnchor(room *liveentity.LiveRoom, anchorId uint64, seatIndex int, lock bool) error {
	if anchorId != room.ID {
		return errercode.CreateCode(errercode.NoPermission)
	}
	if !voiceChatMicIsValidAudienceSeat(seatIndex) {
		return errercode.CreateCode(errercode.VoiceChatMicSeatInvalid)
	}
	changed := false
	mu := voiceChatMicRoomLock(room.ID)
	mu.Lock()
	seat := liveroomdao.GetVoiceChatMicSeat(room.ID, seatIndex)
	if seat == nil || seat.LiveRecordId != room.LiveRecordId {
		mu.Unlock()
		return errercode.CreateCode(errercode.VoiceChatMicSeatInvalid)
	}
	if lock {
		if seat.UserId != 0 {
			mu.Unlock()
			return errercode.CreateCode(errercode.VoiceChatMicSeatTaken)
		}
		if seat.Status != liveentity.VoiceChatMicSeatStatusLocked {
			seat.SetStatus(liveentity.VoiceChatMicSeatStatusLocked)
			liveroomdao.PublishVoiceChatMicSeat(seat)
			changed = true
		}
	} else if liveentity.IsVoiceChatMicSeatLocked(seat) {
		seat.SetStatus(liveentity.VoiceChatMicSeatStatusIdle)
		liveroomdao.PublishVoiceChatMicSeat(seat)
		changed = true
	}
	mu.Unlock()
	if changed {
		broadcastVoiceChatMicState(room.ID)
	}
	return nil
}

func pushVoiceChatMicSeatMuteToUser(roomId, userId uint64, seatIndex int, muted bool) {
	if roomId == 0 || userId == 0 {
		return
	}
	item := &liveroomdto.VoiceChatMicSeatMutePushItem{
		RoomId:    strconv.FormatUint(roomId, 10),
		SeatIndex: seatIndex,
		Muted:     muted,
	}
	push.Data(userId, cmd.LiveRoomVoiceChatMicSeatMute, item)
}

func setVoiceChatMicSeatMutedByAnchor(room *liveentity.LiveRoom, anchorId uint64, seatIndex int, mute bool) error {
	if anchorId != room.ID {
		return errercode.CreateCode(errercode.NoPermission)
	}
	if !voiceChatMicIsValidSeatIndex(seatIndex) {
		return errercode.CreateCode(errercode.VoiceChatMicSeatInvalid)
	}
	changed := false
	var micUserId uint64
	mu := voiceChatMicRoomLock(room.ID)
	mu.Lock()
	seat := liveroomdao.GetVoiceChatMicSeat(room.ID, seatIndex)
	if seat == nil || seat.LiveRecordId != room.LiveRecordId {
		mu.Unlock()
		return errercode.CreateCode(errercode.VoiceChatMicSeatInvalid)
	}
	if seat.UserId == 0 {
		mu.Unlock()
		return errercode.CreateCode(errercode.VoiceChatMicSeatNotOccupied)
	}
	micUserId = seat.UserId
	if mute {
		if seat.Status == liveentity.VoiceChatMicSeatStatusMicMuted {
			mu.Unlock()
			return nil
		}
		if seat.Status != liveentity.VoiceChatMicSeatStatusOnMic {
			mu.Unlock()
			return errercode.CreateCode(errercode.VoiceChatMicSeatNotOccupied)
		}
		seat.SetStatus(liveentity.VoiceChatMicSeatStatusMicMuted)
		liveroomdao.PublishVoiceChatMicSeat(seat)
		changed = true
	} else {
		if seat.Status != liveentity.VoiceChatMicSeatStatusMicMuted {
			mu.Unlock()
			return nil
		}
		seat.SetStatus(liveentity.VoiceChatMicSeatStatusOnMic)
		liveroomdao.PublishVoiceChatMicSeat(seat)
		changed = true
	}
	mu.Unlock()
	if changed {
		pushVoiceChatMicSeatMuteToUser(room.ID, micUserId, seatIndex, mute)
		broadcastVoiceChatMicState(room.ID)
	}
	return nil
}

// MuteVoiceChatMicSeat 主播对指定麦位禁音
func MuteVoiceChatMicSeat(ctx context.Context, req *liveroomdto.MuteVoiceChatMicSeatReq) (*liveroomdto.MuteVoiceChatMicSeatRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	room, err := loadVoiceChatLiveRoom(req.RoomId)
	if err != nil {
		return nil, err
	}
	if err := setVoiceChatMicSeatMutedByAnchor(room, anchorId, req.SeatIndex, true); err != nil {
		return nil, err
	}
	return &liveroomdto.MuteVoiceChatMicSeatRes{Success: true}, nil
}

// UnmuteVoiceChatMicSeat 主播解除麦位禁音
func UnmuteVoiceChatMicSeat(ctx context.Context, req *liveroomdto.UnmuteVoiceChatMicSeatReq) (*liveroomdto.UnmuteVoiceChatMicSeatRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	room, err := loadVoiceChatLiveRoom(req.RoomId)
	if err != nil {
		return nil, err
	}
	if err := setVoiceChatMicSeatMutedByAnchor(room, anchorId, req.SeatIndex, false); err != nil {
		return nil, err
	}
	return &liveroomdto.UnmuteVoiceChatMicSeatRes{Success: true}, nil
}

// InviteVoiceChatMic 主播抱观众上麦
func InviteVoiceChatMic(ctx context.Context, req *liveroomdto.InviteVoiceChatMicReq) (*liveroomdto.InviteVoiceChatMicRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	room, err := loadVoiceChatLiveRoom(req.RoomId)
	if err != nil {
		return nil, err
	}
	if anchorId != room.ID {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	if voiceChatMicModeForRoom(room.ID) == liveentity.VoiceChatMicModeHostOnly {
		return nil, errercode.CreateCode(errercode.VoiceChatMicModeHostOnly)
	}
	if req.UserId == room.ID {
		return nil, errercode.CreateCode(errercode.VoiceChatMicHostSeatLocked)
	}
	if err := requireAudienceOnline(room.ID, req.UserId); err != nil {
		return nil, err
	}
	if _, onSeat := voiceChatMicUserOnSeat(room.ID, req.UserId); onSeat {
		return nil, errercode.CreateCode(errercode.VoiceChatMicAlreadyOn)
	}
	preferredSeat := req.SeatIndex
	if preferredSeat != 0 && !voiceChatMicIsValidAudienceSeat(preferredSeat) {
		return nil, errercode.CreateCode(errercode.VoiceChatMicSeatInvalid)
	}
	seatIndex, ok := voiceChatMicAssignAudienceSeat(room.ID, req.UserId, preferredSeat)
	if !ok {
		return nil, errercode.CreateCode(voiceChatMicAssignAudienceSeatErr(room.ID, room.LiveRecordId, preferredSeat))
	}
	broadcastVoiceChatMicState(room.ID)
	return &liveroomdto.InviteVoiceChatMicRes{Success: true, SeatIndex: seatIndex}, nil
}

// SwitchVoiceChatMicSeat 在麦观众换到另一空麦位
func SwitchVoiceChatMicSeat(ctx context.Context, req *liveroomdto.SwitchVoiceChatMicSeatReq) (*liveroomdto.SwitchVoiceChatMicSeatRes, error) {
	userId := httpserver.GetAuthId(ctx)
	room, err := loadVoiceChatLiveRoom(req.RoomId)
	if err != nil {
		return nil, err
	}
	if userId == room.ID {
		return nil, errercode.CreateCode(errercode.VoiceChatMicHostSeatLocked)
	}
	fromSeat, toSeat, ok := voiceChatMicSwitchAudienceSeat(room.ID, userId, req.SeatIndex)
	if !ok {
		if _, onSeat := voiceChatMicUserOnSeat(room.ID, userId); !onSeat {
			return nil, errercode.CreateCode(errercode.VoiceChatMicNotOnSeat)
		}
		return nil, errercode.CreateCode(voiceChatMicSwitchAudienceSeatErr(room.ID, room.LiveRecordId, req.SeatIndex))
	}
	broadcastVoiceChatMicState(room.ID)
	return &liveroomdto.SwitchVoiceChatMicSeatRes{
		Success:       true,
		FromSeatIndex: fromSeat,
		SeatIndex:     toSeat,
	}, nil
}

func batchSetVoiceChatMicSeatsLockedByAnchor(room *liveentity.LiveRoom, anchorId uint64, seatIndexes []int, lock bool) (int, error) {
	if anchorId != room.ID {
		return 0, errercode.CreateCode(errercode.NoPermission)
	}
	indexes := dedupeAudienceSeatIndexes(seatIndexes)
	if len(indexes) == 0 {
		return 0, errercode.CreateCode(errercode.VoiceChatMicSeatInvalid)
	}
	changed := 0
	mu := voiceChatMicRoomLock(room.ID)
	mu.Lock()
	for _, seatIndex := range indexes {
		seat := liveroomdao.GetVoiceChatMicSeat(room.ID, seatIndex)
		if seat == nil || seat.LiveRecordId != room.LiveRecordId {
			continue
		}
		if lock {
			if seat.UserId != 0 {
				continue
			}
			if seat.Status == liveentity.VoiceChatMicSeatStatusLocked {
				continue
			}
			seat.SetStatus(liveentity.VoiceChatMicSeatStatusLocked)
			liveroomdao.PublishVoiceChatMicSeat(seat)
			changed++
			continue
		}
		if seat.Status != liveentity.VoiceChatMicSeatStatusLocked {
			continue
		}
		seat.SetStatus(liveentity.VoiceChatMicSeatStatusIdle)
		liveroomdao.PublishVoiceChatMicSeat(seat)
		changed++
	}
	mu.Unlock()
	if changed > 0 {
		broadcastVoiceChatMicState(room.ID)
	}
	return changed, nil
}

// BatchLockVoiceChatMicSeats 主播批量锁麦
func BatchLockVoiceChatMicSeats(ctx context.Context, req *liveroomdto.BatchLockVoiceChatMicSeatsReq) (*liveroomdto.BatchLockVoiceChatMicSeatsRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	room, err := loadVoiceChatLiveRoom(req.RoomId)
	if err != nil {
		return nil, err
	}
	count, err := batchSetVoiceChatMicSeatsLockedByAnchor(room, anchorId, req.SeatIndexes, true)
	if err != nil {
		return nil, err
	}
	return &liveroomdto.BatchLockVoiceChatMicSeatsRes{Success: true, Count: count}, nil
}

// BatchUnlockVoiceChatMicSeats 主播批量解锁麦位
func BatchUnlockVoiceChatMicSeats(ctx context.Context, req *liveroomdto.BatchUnlockVoiceChatMicSeatsReq) (*liveroomdto.BatchUnlockVoiceChatMicSeatsRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	room, err := loadVoiceChatLiveRoom(req.RoomId)
	if err != nil {
		return nil, err
	}
	count, err := batchSetVoiceChatMicSeatsLockedByAnchor(room, anchorId, req.SeatIndexes, false)
	if err != nil {
		return nil, err
	}
	return &liveroomdto.BatchUnlockVoiceChatMicSeatsRes{Success: true, Count: count}, nil
}

// ClearVoiceChatMicApplies 主播清空本场全部上麦申请
func ClearVoiceChatMicApplies(ctx context.Context, req *liveroomdto.ClearVoiceChatMicAppliesReq) (*liveroomdto.ClearVoiceChatMicAppliesRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	room, err := loadVoiceChatLiveRoom(req.RoomId)
	if err != nil {
		return nil, err
	}
	if anchorId != room.ID {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	liveroomdao.ClearActiveVoiceChatMicApplies(room.ID, room.LiveRecordId)
	return &liveroomdto.ClearVoiceChatMicAppliesRes{Success: true}, nil
}

// LockVoiceChatMicSeat 主播锁定空观众麦位
func LockVoiceChatMicSeat(ctx context.Context, req *liveroomdto.LockVoiceChatMicSeatReq) (*liveroomdto.LockVoiceChatMicSeatRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	room, err := loadVoiceChatLiveRoom(req.RoomId)
	if err != nil {
		return nil, err
	}
	if err := setVoiceChatMicSeatLockedByAnchor(room, anchorId, req.SeatIndex, true); err != nil {
		return nil, err
	}
	return &liveroomdto.LockVoiceChatMicSeatRes{Success: true}, nil
}

// UnlockVoiceChatMicSeat 主播解锁观众麦位
func UnlockVoiceChatMicSeat(ctx context.Context, req *liveroomdto.UnlockVoiceChatMicSeatReq) (*liveroomdto.UnlockVoiceChatMicSeatRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	room, err := loadVoiceChatLiveRoom(req.RoomId)
	if err != nil {
		return nil, err
	}
	if err := setVoiceChatMicSeatLockedByAnchor(room, anchorId, req.SeatIndex, false); err != nil {
		return nil, err
	}
	return &liveroomdto.UnlockVoiceChatMicSeatRes{Success: true}, nil
}

// RejectVoiceChatMicApply 主播拒绝上麦申请
func RejectVoiceChatMicApply(ctx context.Context, req *liveroomdto.RejectVoiceChatMicApplyReq) (*liveroomdto.RejectVoiceChatMicApplyRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	room, err := loadVoiceChatLiveRoom(req.RoomId)
	if err != nil {
		return nil, err
	}
	if anchorId != room.ID {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	if voiceChatMicModeForRoom(room.ID) != liveentity.VoiceChatMicModeApply {
		return nil, errercode.CreateCode(errercode.VoiceChatMicModeNotApply)
	}

	if !voiceChatMicHasApply(room.ID, req.UserId) {
		return nil, errercode.CreateCode(errercode.VoiceChatMicApplyNotFound)
	}
	liveroomdao.CancelVoiceChatMicApply(room.ID, req.UserId)
	return &liveroomdto.RejectVoiceChatMicApplyRes{Success: true}, nil
}

// LeaveVoiceChatMic 观众主动下麦
func LeaveVoiceChatMic(ctx context.Context, req *liveroomdto.LeaveVoiceChatMicReq) (*liveroomdto.LeaveVoiceChatMicRes, error) {
	userId := httpserver.GetAuthId(ctx)
	room, err := loadVoiceChatLiveRoom(req.RoomId)
	if err != nil {
		return nil, err
	}
	if userId == room.ID {
		return nil, errercode.CreateCode(errercode.VoiceChatMicHostSeatLocked)
	}

	if seatIndex, _ := voiceChatMicUserOnSeat(room.ID, userId); seatIndex == liveentity.VoiceChatHostSeatIndex {
		return nil, errercode.CreateCode(errercode.VoiceChatMicHostSeatLocked)
	}
	if !voiceChatMicClearUserFromSeat(room.ID, room.LiveRecordId, userId) {
		return nil, errercode.CreateCode(errercode.VoiceChatMicNotOnSeat)
	}
	broadcastVoiceChatMicState(room.ID)
	return &liveroomdto.LeaveVoiceChatMicRes{Success: true}, nil
}

// KickVoiceChatMic 主播抱下麦
func KickVoiceChatMic(ctx context.Context, req *liveroomdto.KickVoiceChatMicReq) (*liveroomdto.KickVoiceChatMicRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	room, err := loadVoiceChatLiveRoom(req.RoomId)
	if err != nil {
		return nil, err
	}
	if anchorId != room.ID {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	if req.UserId == room.ID {
		return nil, errercode.CreateCode(errercode.VoiceChatMicHostSeatLocked)
	}

	if seatIndex, _ := voiceChatMicUserOnSeat(room.ID, req.UserId); seatIndex == liveentity.VoiceChatHostSeatIndex {
		return nil, errercode.CreateCode(errercode.VoiceChatMicHostSeatLocked)
	}
	if !voiceChatMicClearUserFromSeat(room.ID, room.LiveRecordId, req.UserId) {
		return nil, errercode.CreateCode(errercode.VoiceChatMicNotOnSeat)
	}
	liveroomdao.CancelVoiceChatMicApply(room.ID, req.UserId)
	broadcastVoiceChatMicState(room.ID)
	return &liveroomdto.KickVoiceChatMicRes{Success: true}, nil
}
