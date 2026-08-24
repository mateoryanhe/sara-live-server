package liveroomdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/live"
)

// RevenueLogCMSListFilter CMS直播收益流水查询条件
type RevenueLogCMSListFilter struct {
	ReceiverId   uint64
	ReceiverIds  []uint64
	LiveRecordId uint64
	RevenueType  uint8
	StartTime    int64
	EndTime      int64
	PageIndex    int
	PageSize     int
}

func (f *RevenueLogCMSListFilter) receiverIds() []uint64 {
	if f == nil {
		return nil
	}
	if len(f.ReceiverIds) > 0 {
		return f.ReceiverIds
	}
	if f.ReceiverId > 0 {
		return []uint64{f.ReceiverId}
	}
	return nil
}

// ParseRevenueLogReceiverIds 合并收益流水查询中的收益用户ID参数(兼容单选与多选)
func ParseRevenueLogReceiverIds(receiverId, platformAnchorId, guildAnchorId string, receiverIds []string) []uint64 {
	return ParseLiveRecordAnchorIds(receiverId, platformAnchorId, guildAnchorId, receiverIds)
}

// RevenueLogCMSList CMS分页查询直播收益流水(按ID倒序)
func RevenueLogCMSList(f *RevenueLogCMSListFilter) (int, []*entity.LiveRevenueLog) {
	list := make([]*entity.LiveRevenueLog, 0)
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
	m := g.Model(string(entity.TbLiveRevenueLog)).Ctx(ctx)
	if receiverIds := f.receiverIds(); len(receiverIds) > 0 {
		m = m.Where(string(entity.LiveRevenueLogReceiverId)+" IN (?)", receiverIds)
	}
	if f.LiveRecordId > 0 {
		m = m.Where(string(entity.LiveRevenueLogLiveRecordId)+" = ?", f.LiveRecordId)
	}
	if f.RevenueType > 0 {
		m = m.Where(string(entity.LiveRevenueLogRevenueType)+" = ?", f.RevenueType)
	}
	if f.StartTime > 0 {
		m = m.Where("created_at >= ?", time.Unix(f.StartTime, 0))
	}
	if f.EndTime > 0 {
		m = m.Where("created_at <= ?", time.Unix(f.EndTime, 0))
	}
	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	_ = m.Clone().Order("id desc").
		Limit(f.PageSize).Offset((f.PageIndex - 1) * f.PageSize).
		Scan(&list)
	return total, list
}
