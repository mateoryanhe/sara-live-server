package liveroomdao

import (
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/entity/live"
)

var (
	incomeUnsettledCache = gmap.NewKVMap[uint64, *entity.LiveRoomIncomeUnsettled](false)
	incomeSettledCache   = gmap.NewKVMap[uint64, *entity.LiveRoomIncomeSettled](false)
	incomeTotalCache     = gmap.NewKVMap[uint64, *entity.LiveRoomIncomeTotal](false)
)

func initLiveRoomIncomeDao() {
	// 随直播间预热:在 initLiveRoomDao 加载上架房间后调用 PreloadLiveRoomIncomes
}

// PreloadLiveRoomIncomes 预热上架直播间的三张收益表(永久缓存)
func PreloadLiveRoomIncomes(roomIds []uint64) {
	if len(roomIds) == 0 {
		return
	}
	unsettledRows := make([]*entity.LiveRoomIncomeUnsettled, 0, len(roomIds))
	_ = g.Model(string(entity.TbLiveRoomIncomeUnsettled)).Unscoped().
		WhereIn(string(db.IdName), roomIds).Scan(&unsettledRows)
	settledRows := make([]*entity.LiveRoomIncomeSettled, 0, len(roomIds))
	_ = g.Model(string(entity.TbLiveRoomIncomeSettled)).Unscoped().
		WhereIn(string(db.IdName), roomIds).Scan(&settledRows)
	totalRows := make([]*entity.LiveRoomIncomeTotal, 0, len(roomIds))
	_ = g.Model(string(entity.TbLiveRoomIncomeTotal)).Unscoped().
		WhereIn(string(db.IdName), roomIds).Scan(&totalRows)

	unsettledMap := make(map[uint64]*entity.LiveRoomIncomeUnsettled, len(unsettledRows))
	for _, row := range unsettledRows {
		if row != nil && row.ID != 0 {
			unsettledMap[row.ID] = row
		}
	}
	settledMap := make(map[uint64]*entity.LiveRoomIncomeSettled, len(settledRows))
	for _, row := range settledRows {
		if row != nil && row.ID != 0 {
			settledMap[row.ID] = row
		}
	}
	totalMap := make(map[uint64]*entity.LiveRoomIncomeTotal, len(totalRows))
	for _, row := range totalRows {
		if row != nil && row.ID != 0 {
			totalMap[row.ID] = row
		}
	}

	for _, id := range roomIds {
		if id == 0 {
			continue
		}
		if row, ok := unsettledMap[id]; ok {
			incomeUnsettledCache.Set(id, row)
		} else {
			empty := &entity.LiveRoomIncomeUnsettled{}
			empty.ID = id
			incomeUnsettledCache.Set(id, empty)
		}
		if row, ok := settledMap[id]; ok {
			incomeSettledCache.Set(id, row)
		} else {
			empty := &entity.LiveRoomIncomeSettled{}
			empty.ID = id
			incomeSettledCache.Set(id, empty)
		}
		if row, ok := totalMap[id]; ok {
			incomeTotalCache.Set(id, row)
		} else {
			empty := &entity.LiveRoomIncomeTotal{}
			empty.ID = id
			incomeTotalCache.Set(id, empty)
		}
	}
}

// AddLiveRoomIncomeToCache 上架时加载/创建三张收益缓存
func AddLiveRoomIncomeToCache(roomId uint64) {
	if roomId == 0 {
		return
	}
	PreloadLiveRoomIncomes([]uint64{roomId})
}

// RemoveLiveRoomIncomeFromCache 下架时移除三张收益缓存
func RemoveLiveRoomIncomeFromCache(roomId uint64) {
	if roomId == 0 {
		return
	}
	incomeUnsettledCache.Remove(roomId)
	incomeSettledCache.Remove(roomId)
	incomeTotalCache.Remove(roomId)
}

// GetLiveRoomIncomeUnsettled 未结算收益(仅内存,无则新建并写入缓存)
func GetLiveRoomIncomeUnsettled(roomId uint64) *entity.LiveRoomIncomeUnsettled {
	if roomId == 0 {
		return nil
	}
	if incomeUnsettledCache.Contains(roomId) {
		return incomeUnsettledCache.Get(roomId)
	}
	row := entity.NewLiveRoomIncomeUnsettled(roomId)
	incomeUnsettledCache.Set(roomId, row)
	return row
}

// GetLiveRoomIncomeSettled 已结算收益(仅内存,无则新建并写入缓存)
func GetLiveRoomIncomeSettled(roomId uint64) *entity.LiveRoomIncomeSettled {
	if roomId == 0 {
		return nil
	}
	if incomeSettledCache.Contains(roomId) {
		return incomeSettledCache.Get(roomId)
	}
	row := entity.NewLiveRoomIncomeSettled(roomId)
	incomeSettledCache.Set(roomId, row)
	return row
}

// GetLiveRoomIncomeTotal 生涯累计收益(仅内存,无则新建并写入缓存)
func GetLiveRoomIncomeTotal(roomId uint64) *entity.LiveRoomIncomeTotal {
	if roomId == 0 {
		return nil
	}
	if incomeTotalCache.Contains(roomId) {
		return incomeTotalCache.Get(roomId)
	}
	row := entity.NewLiveRoomIncomeTotal(roomId)
	incomeTotalCache.Set(roomId, row)
	return row
}

// GetLiveRoomIncomeTotalFromCache 仅读缓存,未命中返回 nil(不查库、不新建)
func GetLiveRoomIncomeTotalFromCache(roomId uint64) *entity.LiveRoomIncomeTotal {
	if roomId == 0 || !incomeTotalCache.Contains(roomId) {
		return nil
	}
	return incomeTotalCache.Get(roomId)
}

// GetLiveRoomIncomeTotalFromDB 直查数据库(回收站等场景)
func GetLiveRoomIncomeTotalFromDB(roomId uint64) *entity.LiveRoomIncomeTotal {
	if roomId == 0 {
		return nil
	}
	var row entity.LiveRoomIncomeTotal
	err := g.Model(string(entity.TbLiveRoomIncomeTotal)).Unscoped().WherePri(roomId).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// ListLiveRoomIncomeTotalFromDB 批量直查生涯累计收益(回收站)
func ListLiveRoomIncomeTotalFromDB(roomIds []uint64) map[uint64]*entity.LiveRoomIncomeTotal {
	ret := make(map[uint64]*entity.LiveRoomIncomeTotal, len(roomIds))
	if len(roomIds) == 0 {
		return ret
	}
	rows := make([]*entity.LiveRoomIncomeTotal, 0, len(roomIds))
	_ = g.Model(string(entity.TbLiveRoomIncomeTotal)).Unscoped().
		WhereIn(string(db.IdName), roomIds).Scan(&rows)
	for _, row := range rows {
		if row != nil && row.ID != 0 {
			ret[row.ID] = row
		}
	}
	return ret
}

// GetLiveRoomIncomeUnsettledForCMS 未结算收益(缓存优先,否则直查DB,不新建)
func GetLiveRoomIncomeUnsettledForCMS(roomId uint64) *entity.LiveRoomIncomeUnsettled {
	if roomId == 0 {
		return nil
	}
	if incomeUnsettledCache.Contains(roomId) {
		return incomeUnsettledCache.Get(roomId)
	}
	var row entity.LiveRoomIncomeUnsettled
	err := g.Model(string(entity.TbLiveRoomIncomeUnsettled)).Unscoped().WherePri(roomId).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// ListLiveRoomIncomeUnsettledTotalForCMS 批量查询未结算总收益(缓存优先,否则直查DB,不新建)
func ListLiveRoomIncomeUnsettledTotalForCMS(roomIds []uint64) map[uint64]float64 {
	ret := make(map[uint64]float64)
	if len(roomIds) == 0 {
		return ret
	}
	seen := make(map[uint64]struct{}, len(roomIds))
	missing := make([]uint64, 0, len(roomIds))
	for _, roomId := range roomIds {
		if roomId == 0 {
			continue
		}
		if _, ok := seen[roomId]; ok {
			continue
		}
		seen[roomId] = struct{}{}
		if incomeUnsettledCache.Contains(roomId) {
			if row := incomeUnsettledCache.Get(roomId); row != nil {
				ret[roomId] = row.TotalIncome
			}
			continue
		}
		missing = append(missing, roomId)
	}
	if len(missing) == 0 {
		return ret
	}
	rows := make([]*entity.LiveRoomIncomeUnsettled, 0, len(missing))
	_ = g.Model(string(entity.TbLiveRoomIncomeUnsettled)).Unscoped().
		WhereIn(string(db.IdName), missing).Scan(&rows)
	for _, row := range rows {
		if row != nil && row.ID != 0 {
			ret[row.ID] = row.TotalIncome
		}
	}
	return ret
}

// GetLiveRoomIncomeSettledForCMS 已结算收益(缓存优先,否则直查DB,不新建)
func GetLiveRoomIncomeSettledForCMS(roomId uint64) *entity.LiveRoomIncomeSettled {
	if roomId == 0 {
		return nil
	}
	if incomeSettledCache.Contains(roomId) {
		return incomeSettledCache.Get(roomId)
	}
	var row entity.LiveRoomIncomeSettled
	err := g.Model(string(entity.TbLiveRoomIncomeSettled)).Unscoped().WherePri(roomId).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// GetLiveRoomIncomeTotalForCMS 生涯累计收益(缓存优先,否则直查DB,不新建)
func GetLiveRoomIncomeTotalForCMS(roomId uint64) *entity.LiveRoomIncomeTotal {
	if roomId == 0 {
		return nil
	}
	if incomeTotalCache.Contains(roomId) {
		return incomeTotalCache.Get(roomId)
	}
	return GetLiveRoomIncomeTotalFromDB(roomId)
}

// ListLiveRoomIncomeUnsettledArchives 查询直播间下架未结算归档记录
func ListLiveRoomIncomeUnsettledArchives(roomId uint64, limit int) []*entity.LiveRoomIncomeUnsettledArchive {
	if roomId == 0 {
		return nil
	}
	if limit <= 0 {
		limit = 50
	}
	rows := make([]*entity.LiveRoomIncomeUnsettledArchive, 0)
	_ = g.Model(string(entity.TbLiveRoomIncomeUnsettledArchive)).
		Where(string(entity.LiveRoomIncomeUnsettledArchiveRoomId)+" = ?", roomId).
		Order(string(db.CreatedAtName) + " desc").
		Limit(limit).
		Scan(&rows)
	return rows
}

// getLiveRoomIncomeUnsettledForArchive 下架归档用:优先缓存,否则直查DB(不新建)
func getLiveRoomIncomeUnsettledForArchive(roomId uint64) *entity.LiveRoomIncomeUnsettled {
	if roomId == 0 {
		return nil
	}
	if incomeUnsettledCache.Contains(roomId) {
		return incomeUnsettledCache.Get(roomId)
	}
	var row entity.LiveRoomIncomeUnsettled
	err := g.Model(string(entity.TbLiveRoomIncomeUnsettled)).Unscoped().WherePri(roomId).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// ArchiveAndClearUnsettledIncome 下架时:新建一条归档记录并清零未结算表
func ArchiveAndClearUnsettledIncome(roomId, guildId uint64) {
	unsettled := getLiveRoomIncomeUnsettledForArchive(roomId)
	if unsettled == nil || unsettled.IsZero() {
		return
	}
	snap := unsettled.SnapshotAndClear()
	if snap.IsZero() {
		return
	}
	_ = entity.NewLiveRoomIncomeUnsettledArchive(roomId, guildId, &snap)
}
