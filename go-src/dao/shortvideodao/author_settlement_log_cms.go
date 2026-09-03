package shortvideodao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/shortvideo"
)

// AuthorSettlementLogCMSListFilter CMS 短视频作者结算日志查询条件
type AuthorSettlementLogCMSListFilter struct {
	UserId    uint64
	StartTime int64
	EndTime   int64
	PageIndex int
	PageSize  int
}

// AuthorSettlementLogCMSList CMS 分页查询短视频作者结算日志（按 ID 倒序）
func AuthorSettlementLogCMSList(f *AuthorSettlementLogCMSListFilter) (int, []*entity.ShortVideoAuthorSettlementLog) {
	list := make([]*entity.ShortVideoAuthorSettlementLog, 0)
	if f == nil {
		return 0, list
	}
	pageIndex := f.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	pageSize := f.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	ctx := gctx.New()
	m := g.Model(string(entity.TbShortVideoAuthorSettlementLog)).Ctx(ctx)
	if f.UserId > 0 {
		m = m.Where(string(entity.ShortVideoAuthorSettlementLogUserId)+" = ?", f.UserId)
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
		Limit(pageSize).Offset((pageIndex - 1) * pageSize).
		Scan(&list)
	return total, list
}
