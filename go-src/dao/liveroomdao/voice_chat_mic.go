package liveroomdao

import (
	"context"
	"sort"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/cache"
	liveentity "xr-game-server/entity/live"
)

var (
	voiceChatMicSeatCacheMgr     *cache.RowCache[*liveentity.VoiceChatMicSeat]
	voiceChatMicApplyRowCacheMgr *cache.RowCache[*liveentity.VoiceChatMicApply]
	voiceChatMicApplyCacheMgr    *cache.ListCache[*liveentity.VoiceChatMicApply]
)

func initLiveRoomVoiceChatMicDao() {
	voiceChatMicSeatCacheMgr = cache.NewRowCache[*liveentity.VoiceChatMicSeat]()
	voiceChatMicApplyRowCacheMgr = cache.NewRowCache[*liveentity.VoiceChatMicApply]()
	voiceChatMicApplyCacheMgr = cache.NewPermanentListCache[*liveentity.VoiceChatMicApply]()
}

func isVoiceChatRoomCfg(roomId uint64) bool {
	cfg := GetLiveRoomCfgFromCache(roomId)
	if cfg == nil {
		cfg = GetLiveRoomCfg(roomId)
	}
	return cfg != nil && cfg.Category == liveentity.LiveRoomCategoryVoiceChat
}

// PreloadVoiceChatMicLiveSessionsFromDB 热重启/冷启动:恢复进行中语聊场麦位与申请(不触发重新开播).
func PreloadVoiceChatMicLiveSessionsFromDB() int {
	restored := 0
	for _, room := range GetAllLiveRoom() {
		if room == nil || room.ID == 0 || room.LiveRecordId == 0 {
			continue
		}
		if !isVoiceChatRoomCfg(room.ID) {
			continue
		}
		RestoreVoiceChatMicSessionFromDB(room.ID, room.LiveRecordId, room.ID)
		restored++
	}
	return restored
}

func PublishVoiceChatMicSeat(row *liveentity.VoiceChatMicSeat) {
	if row == nil || row.ID == "" || voiceChatMicSeatCacheMgr == nil {
		return
	}
	voiceChatMicSeatCacheMgr.PublishRow(gctx.New(), row.ID, row)
}

func PublishVoiceChatMicApply(row *liveentity.VoiceChatMicApply) {
	if row == nil || row.ID == "" || voiceChatMicApplyRowCacheMgr == nil {
		return
	}
	voiceChatMicApplyRowCacheMgr.PublishRow(gctx.New(), row.ID, row)
}

func placeholderVoiceChatMicSeat(roomId uint64, seatIndex int) *liveentity.VoiceChatMicSeat {
	return &liveentity.VoiceChatMicSeat{
		ID:        liveentity.BuildVoiceChatMicSeatId(roomId, seatIndex),
		RoomId:    roomId,
		SeatIndex: uint8(seatIndex),
		Status:    liveentity.VoiceChatMicSeatStatusIdle,
	}
}

func GetVoiceChatMicSeat(roomId uint64, seatIndex int) *liveentity.VoiceChatMicSeat {
	if roomId == 0 || seatIndex < 0 || seatIndex >= liveentity.VoiceChatSeatCount {
		return nil
	}
	id := liveentity.BuildVoiceChatMicSeatId(roomId, seatIndex)
	if voiceChatMicSeatCacheMgr == nil {
		return loadVoiceChatMicSeatFromDB(id)
	}
	return voiceChatMicSeatCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*liveentity.VoiceChatMicSeat, error) {
		if loaded := loadVoiceChatMicSeatFromDB(id); loaded != nil {
			return loaded, nil
		}
		// 未落库时仅返回内存占位,首次 Set* 时经 syndb 入库
		return placeholderVoiceChatMicSeat(roomId, seatIndex), nil
	})
}

func loadVoiceChatMicSeatsBySession(roomId, liveRecordId uint64) []*liveentity.VoiceChatMicSeat {
	rows := make([]*liveentity.VoiceChatMicSeat, 0)
	if roomId == 0 || liveRecordId == 0 {
		return rows
	}
	_ = g.Model(string(liveentity.TbVoiceChatMicSeat)).
		Where("room_id = ? AND live_record_id = ?", roomId, liveRecordId).
		Order("seat_index asc").
		Scan(&rows)
	return rows
}

func publishVoiceChatMicSeatSnapshot(row *liveentity.VoiceChatMicSeat) {
	if row == nil || row.ID == "" {
		return
	}
	PublishVoiceChatMicSeat(row)
}

// RestoreVoiceChatMicSessionFromDB 将指定场次麦位/申请灌入内存缓存(DB 为准,不取消申请).
func RestoreVoiceChatMicSessionFromDB(roomId, liveRecordId, anchorId uint64) {
	if roomId == 0 || liveRecordId == 0 {
		return
	}
	dbSeats := loadVoiceChatMicSeatsBySession(roomId, liveRecordId)
	byIndex := make(map[int]*liveentity.VoiceChatMicSeat, len(dbSeats))
	for _, seat := range dbSeats {
		if seat == nil {
			continue
		}
		byIndex[int(seat.SeatIndex)] = seat
	}

	if len(dbSeats) == 0 {
		bootstrapVoiceChatMicSessionRows(roomId, liveRecordId, anchorId)
	} else {
		for i := 0; i < liveentity.VoiceChatSeatCount; i++ {
			if seat := byIndex[i]; seat != nil {
				publishVoiceChatMicSeatSnapshot(seat)
				continue
			}
			seat := placeholderVoiceChatMicSeat(roomId, i)
			seat.SetLiveRecordId(liveRecordId)
			seat.SetUserId(0)
			seat.SetStatus(liveentity.VoiceChatMicSeatStatusIdle)
			publishVoiceChatMicSeatSnapshot(seat)
		}
		repairVoiceChatHostMicSeat(roomId, liveRecordId, anchorId)
	}

	list := loadActiveVoiceChatMicAppliesFromDB(roomId, liveRecordId)
	for _, row := range list {
		PublishVoiceChatMicApply(row)
	}
	refreshVoiceChatMicApplyListCacheFromDB(roomId, liveRecordId)
}

func bootstrapVoiceChatMicSessionRows(roomId, liveRecordId, anchorId uint64) {
	for i := 0; i < liveentity.VoiceChatSeatCount; i++ {
		userId := uint64(0)
		if i == liveentity.VoiceChatHostSeatIndex {
			userId = anchorId
		}
		id := liveentity.BuildVoiceChatMicSeatId(roomId, i)
		row := loadVoiceChatMicSeatFromDB(id)
		if row != nil && row.LiveRecordId == liveRecordId {
			publishVoiceChatMicSeatSnapshot(row)
			continue
		}
		row = liveentity.NewVoiceChatMicSeat(roomId, i, liveRecordId, userId)
		publishVoiceChatMicSeatSnapshot(row)
	}
}

func repairVoiceChatHostMicSeat(roomId, liveRecordId, anchorId uint64) {
	if anchorId == 0 {
		return
	}
	seat := GetVoiceChatMicSeat(roomId, liveentity.VoiceChatHostSeatIndex)
	if seat == nil || seat.LiveRecordId != liveRecordId {
		return
	}
	if seat.UserId == anchorId {
		return
	}
	if seat.UserId != 0 && seat.UserId != anchorId {
		return
	}
	seat.SetUserId(anchorId)
	seat.SetStatus(liveentity.VoiceChatMicSeatStatusOnMic)
	publishVoiceChatMicSeatSnapshot(seat)
}

func loadVoiceChatMicSeatFromDB(id string) *liveentity.VoiceChatMicSeat {
	if id == "" {
		return nil
	}
	var row liveentity.VoiceChatMicSeat
	if err := g.Model(string(liveentity.TbVoiceChatMicSeat)).Where("id = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == "" {
		return nil
	}
	return &row
}

func loadActiveVoiceChatMicAppliesFromDB(roomId, liveRecordId uint64) []*liveentity.VoiceChatMicApply {
	rows := make([]*liveentity.VoiceChatMicApply, 0)
	if roomId == 0 || liveRecordId == 0 {
		return rows
	}
	_ = g.Model(string(liveentity.TbVoiceChatMicApply)).
		Where("room_id = ? AND live_record_id = ? AND status = ?",
			roomId, liveRecordId, liveentity.VoiceChatMicApplyStatusActive).
		Order("applied_at asc").
		Scan(&rows)
	return rows
}

// GetActiveVoiceChatMicApplies 查询本场有效上麦申请:优先 ListCache,未命中再查 DB.
func GetActiveVoiceChatMicApplies(roomId, liveRecordId uint64) []*liveentity.VoiceChatMicApply {
	if roomId == 0 || liveRecordId == 0 {
		return make([]*liveentity.VoiceChatMicApply, 0)
	}
	cacheKey := voiceChatMicApplyListCacheKey(roomId, liveRecordId)
	if voiceChatMicApplyCacheMgr == nil {
		return loadActiveVoiceChatMicAppliesFromDB(roomId, liveRecordId)
	}
	list := voiceChatMicApplyCacheMgr.MustGetList(gctx.New(), cacheKey, func(ctx context.Context) ([]*liveentity.VoiceChatMicApply, error) {
		return loadActiveVoiceChatMicAppliesFromDB(roomId, liveRecordId), nil
	})
	if list == nil {
		return make([]*liveentity.VoiceChatMicApply, 0)
	}
	return list
}

func publishVoiceChatMicApplyList(roomId, liveRecordId uint64, list []*liveentity.VoiceChatMicApply) {
	if voiceChatMicApplyCacheMgr == nil || roomId == 0 || liveRecordId == 0 {
		return
	}
	if list == nil {
		list = make([]*liveentity.VoiceChatMicApply, 0)
	}
	voiceChatMicApplyCacheMgr.PublishList(gctx.New(), voiceChatMicApplyListCacheKey(roomId, liveRecordId), list)
}

func cloneVoiceChatMicApplyList(list []*liveentity.VoiceChatMicApply) []*liveentity.VoiceChatMicApply {
	if len(list) == 0 {
		return make([]*liveentity.VoiceChatMicApply, 0)
	}
	out := make([]*liveentity.VoiceChatMicApply, 0, len(list))
	for _, row := range list {
		if row == nil {
			continue
		}
		if row.Status != liveentity.VoiceChatMicApplyStatusActive {
			continue
		}
		out = append(out, row)
	}
	return out
}

func sortVoiceChatMicApplyList(list []*liveentity.VoiceChatMicApply) {
	sort.Slice(list, func(i, j int) bool {
		a, b := list[i], list[j]
		if a == nil || a.AppliedAt == nil {
			return false
		}
		if b == nil || b.AppliedAt == nil {
			return true
		}
		return a.AppliedAt.Before(*b.AppliedAt)
	})
}

func mergeActiveVoiceChatMicApply(list []*liveentity.VoiceChatMicApply, row *liveentity.VoiceChatMicApply, liveRecordId uint64) []*liveentity.VoiceChatMicApply {
	out := cloneVoiceChatMicApplyList(list)
	if row == nil || row.Status != liveentity.VoiceChatMicApplyStatusActive || row.LiveRecordId != liveRecordId {
		sortVoiceChatMicApplyList(out)
		return out
	}
	replaced := false
	for i, item := range out {
		if item != nil && item.UserId == row.UserId {
			out[i] = row
			replaced = true
			break
		}
	}
	if !replaced {
		out = append(out, row)
	}
	sortVoiceChatMicApplyList(out)
	return out
}

func removeActiveVoiceChatMicApply(list []*liveentity.VoiceChatMicApply, userId uint64) []*liveentity.VoiceChatMicApply {
	if userId == 0 {
		return cloneVoiceChatMicApplyList(list)
	}
	out := make([]*liveentity.VoiceChatMicApply, 0, len(list))
	for _, item := range list {
		if item == nil || item.UserId == userId {
			continue
		}
		if item.Status != liveentity.VoiceChatMicApplyStatusActive {
			continue
		}
		out = append(out, item)
	}
	return out
}

// ActiveVoiceChatMicApplyCount 本场有效申请数量(走申请列表缓存).
func ActiveVoiceChatMicApplyCount(roomId, liveRecordId uint64) int {
	return len(GetActiveVoiceChatMicApplies(roomId, liveRecordId))
}

// HasActiveVoiceChatMicApply 是否已有有效申请(走申请列表缓存).
func HasActiveVoiceChatMicApply(roomId, userId, liveRecordId uint64) bool {
	if userId == 0 || liveRecordId == 0 {
		return false
	}
	for _, row := range GetActiveVoiceChatMicApplies(roomId, liveRecordId) {
		if row != nil && row.UserId == userId {
			return true
		}
	}
	return false
}

func voiceChatMicApplyListCacheKey(roomId, liveRecordId uint64) uint64 {
	return roomId*1_000_000_000_000 + liveRecordId
}

// refreshVoiceChatMicApplyListCacheFromDB 仅用于启动恢复等需要从 DB 重建列表的场景.
func refreshVoiceChatMicApplyListCacheFromDB(roomId, liveRecordId uint64) {
	if roomId == 0 || liveRecordId == 0 {
		return
	}
	publishVoiceChatMicApplyList(roomId, liveRecordId, loadActiveVoiceChatMicAppliesFromDB(roomId, liveRecordId))
}

// InitVoiceChatMicSession 新开播时初始化麦位(清空观众位并取消旧申请);热重启恢复请用 RestoreVoiceChatMicSessionFromDB.
func InitVoiceChatMicSession(roomId, liveRecordId, anchorId uint64) {
	if roomId == 0 || liveRecordId == 0 {
		return
	}
	for i := 0; i < liveentity.VoiceChatSeatCount; i++ {
		userId := uint64(0)
		if i == liveentity.VoiceChatHostSeatIndex {
			userId = anchorId
		}
		id := liveentity.BuildVoiceChatMicSeatId(roomId, i)
		row := loadVoiceChatMicSeatFromDB(id)
		if row == nil {
			row = liveentity.NewVoiceChatMicSeat(roomId, i, liveRecordId, userId)
		} else {
			row.SetLiveRecordId(liveRecordId)
			row.SetUserId(userId)
			if i != liveentity.VoiceChatHostSeatIndex {
				row.SetStatus(liveentity.VoiceChatMicSeatStatusIdle)
			} else {
				row.SetStatus(liveentity.VoiceChatMicSeatStatusOnMic)
			}
		}
		publishVoiceChatMicSeatSnapshot(row)
	}
	cancelAllVoiceChatMicApplies(roomId, liveRecordId)
	publishVoiceChatMicApplyList(roomId, liveRecordId, make([]*liveentity.VoiceChatMicApply, 0))
}

// ClearVoiceChatMicSession 下播时清空麦位与申请
func ClearVoiceChatMicSession(roomId, liveRecordId uint64) {
	if roomId == 0 {
		return
	}
	for i := 0; i < liveentity.VoiceChatSeatCount; i++ {
		seat := GetVoiceChatMicSeat(roomId, i)
		if seat == nil {
			continue
		}
		if liveRecordId > 0 && seat.LiveRecordId != liveRecordId {
			continue
		}
		seat.SetUserId(0)
		seat.SetLiveRecordId(0)
		seat.SetStatus(liveentity.VoiceChatMicSeatStatusIdle)
		PublishVoiceChatMicSeat(seat)
	}
	if liveRecordId > 0 {
		cancelAllVoiceChatMicApplies(roomId, liveRecordId)
	}
}

// ClearActiveVoiceChatMicApplies 取消本场全部有效上麦申请
func ClearActiveVoiceChatMicApplies(roomId, liveRecordId uint64) {
	cancelAllVoiceChatMicApplies(roomId, liveRecordId)
}

func cancelAllVoiceChatMicApplies(roomId, liveRecordId uint64) {
	for _, row := range GetActiveVoiceChatMicApplies(roomId, liveRecordId) {
		if row == nil {
			continue
		}
		row.SetStatus(liveentity.VoiceChatMicApplyStatusCancelled)
		PublishVoiceChatMicApply(row)
	}
	publishVoiceChatMicApplyList(roomId, liveRecordId, make([]*liveentity.VoiceChatMicApply, 0))
}

func UpsertActiveVoiceChatMicApply(roomId, userId, liveRecordId uint64, appliedAt time.Time) *liveentity.VoiceChatMicApply {
	row := ResolveVoiceChatMicApplyRow(roomId, userId)
	if row == nil || row.ID == "" {
		row = liveentity.NewVoiceChatMicApply(roomId, userId, liveRecordId, appliedAt)
	} else {
		row.SetLiveRecordId(liveRecordId)
		row.SetStatus(liveentity.VoiceChatMicApplyStatusActive)
		row.SetAppliedAt(&appliedAt)
	}
	PublishVoiceChatMicApply(row)
	list := mergeActiveVoiceChatMicApply(GetActiveVoiceChatMicApplies(roomId, liveRecordId), row, liveRecordId)
	publishVoiceChatMicApplyList(roomId, liveRecordId, list)
	return row
}

func CancelVoiceChatMicApply(roomId, userId uint64) {
	room := GetRoomById(roomId)
	if room == nil || room.LiveRecordId == 0 {
		return
	}
	liveRecordId := room.LiveRecordId
	row := ResolveVoiceChatMicApplyRow(roomId, userId)
	if row == nil || row.Status != liveentity.VoiceChatMicApplyStatusActive || row.LiveRecordId != liveRecordId {
		return
	}
	row.SetStatus(liveentity.VoiceChatMicApplyStatusCancelled)
	PublishVoiceChatMicApply(row)
	list := removeActiveVoiceChatMicApply(GetActiveVoiceChatMicApplies(roomId, liveRecordId), userId)
	publishVoiceChatMicApplyList(roomId, liveRecordId, list)
}

// ResolveVoiceChatMicApplyRow 单条申请:优先 RowCache,未命中再查 DB.
func ResolveVoiceChatMicApplyRow(roomId, userId uint64) *liveentity.VoiceChatMicApply {
	if roomId == 0 || userId == 0 {
		return nil
	}
	id := liveentity.BuildVoiceChatMicApplyId(roomId, userId)
	if voiceChatMicApplyRowCacheMgr == nil {
		return loadVoiceChatMicApplyFromDB(id)
	}
	return voiceChatMicApplyRowCacheMgr.MustGetRow(gctx.New(), id, func(ctx context.Context) (*liveentity.VoiceChatMicApply, error) {
		if loaded := loadVoiceChatMicApplyFromDB(id); loaded != nil {
			return loaded, nil
		}
		return nil, nil
	})
}

func GetVoiceChatMicApply(roomId, userId uint64) *liveentity.VoiceChatMicApply {
	row := ResolveVoiceChatMicApplyRow(roomId, userId)
	if row == nil || row.ID == "" {
		return nil
	}
	return row
}

func loadVoiceChatMicApplyFromDB(id string) *liveentity.VoiceChatMicApply {
	if id == "" {
		return nil
	}
	var row liveentity.VoiceChatMicApply
	if err := g.Model(string(liveentity.TbVoiceChatMicApply)).Where("id = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == "" {
		return nil
	}
	return &row
}
