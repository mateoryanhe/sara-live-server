package liveroomdao

import (
	"strconv"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/live"
)

// LiveRecordCMSListFilter CMS直播记录查询条件
type LiveRecordCMSListFilter struct {
	AnchorIds []uint64
	StartTime int64
	EndTime   int64
	PageIndex int
	PageSize  int
}

func (f *LiveRecordCMSListFilter) anchorIds() []uint64 {
	if f == nil {
		return nil
	}
	return f.AnchorIds
}

func parseUint64AnchorId(val string) uint64 {
	if val == "" {
		return 0
	}
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return id
}

// ParseLiveRecordAnchorIds 合并直播记录查询中的主播ID参数(兼容单选与多选)
func ParseLiveRecordAnchorIds(anchorId, platformAnchorId, guildAnchorId string, anchorIds []string) []uint64 {
	seen := make(map[uint64]struct{})
	ret := make([]uint64, 0, len(anchorIds)+3)
	add := func(id uint64) {
		if id == 0 {
			return
		}
		if _, ok := seen[id]; ok {
			return
		}
		seen[id] = struct{}{}
		ret = append(ret, id)
	}
	add(parseUint64AnchorId(anchorId))
	add(parseUint64AnchorId(platformAnchorId))
	add(parseUint64AnchorId(guildAnchorId))
	for _, val := range anchorIds {
		add(parseUint64AnchorId(val))
	}
	return ret
}

func mergeLiveRecordsFromCache(rows []*entity.LiveRecord) []*entity.LiveRecord {
	if len(rows) == 0 {
		return rows
	}
	list := make([]*entity.LiveRecord, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		list = append(list, resolveLiveRecordTarget(row))
	}
	return list
}

// LiveRecordCMSList CMS分页查询直播记录(按ID倒序;DB结果逐条合并单条缓存)
func LiveRecordCMSList(f *LiveRecordCMSListFilter) (int, []*entity.LiveRecord) {
	list := make([]*entity.LiveRecord, 0)
	if f == nil {
		return 0, list
	}
	if f.PageIndex <= 0 {
		f.PageIndex = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 20
	}
	ctx := gctx.New()
	m := g.Model(string(entity.TbLiveRecord)).Ctx(ctx)
	if anchorIds := f.anchorIds(); len(anchorIds) > 0 {
		m = m.Where(string(entity.LiveRecordAnchorId)+" IN (?)", anchorIds)
	}
	if f.StartTime > 0 {
		m = m.Where(string(entity.LiveRecordStartTime)+" >= ?", time.Unix(f.StartTime, 0))
	}
	if f.EndTime > 0 {
		m = m.Where(string(entity.LiveRecordStartTime)+" <= ?", time.Unix(f.EndTime, 0))
	}
	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	_ = m.Clone().Order("id desc").
		Limit(f.PageSize).Offset((f.PageIndex - 1) * f.PageSize).
		Scan(&list)
	return total, mergeLiveRecordsFromCache(list)
}
