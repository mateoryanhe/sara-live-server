package liveroomdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity"
)

// GiftLogCMSListFilter CMS礼物流水查询条件
type GiftLogCMSListFilter struct {
	ReceiverId uint64
	StartTime  int64
	EndTime    int64
	PageIndex  int
	PageSize   int
}

// GiftLogCMSList CMS分页查询礼物流水(按ID倒序)
func GiftLogCMSList(f *GiftLogCMSListFilter) (int, []*entity.LiveGiftLog) {
	list := make([]*entity.LiveGiftLog, 0)
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
	m := g.Model(string(entity.TbLiveGiftLog)).Ctx(ctx)
	if f.ReceiverId > 0 {
		m = m.Where(string(entity.LiveGiftLogReceiverId)+" = ?", f.ReceiverId)
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
