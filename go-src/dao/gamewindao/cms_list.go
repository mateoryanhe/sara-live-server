package gamewindao

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity"
)

// CMSListFilter CMS 派彩记录查询条件
type CMSListFilter struct {
	UserId       uint64
	GameCode     string
	OrderId      string
	PlatformType string
	PageIndex    int
	PageSize     int
}

// CMSList CMS 分页查询游戏派彩记录(按 ID 倒序)
func CMSList(f *CMSListFilter) (int, []*entity.GameWinLog) {
	list := make([]*entity.GameWinLog, 0)
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
	m := g.Model(string(entity.TbGameWinLog)).Ctx(ctx).
		Where(string(entity.GameWinLogAmount)+" > ?", 0)
	if f.UserId > 0 {
		m = m.Where(string(entity.GameWinLogUserId)+" = ?", f.UserId)
	}
	if gameCode := strings.TrimSpace(f.GameCode); gameCode != "" {
		m = m.Where(string(entity.GameWinLogGameCode)+" LIKE ?", "%"+gameCode+"%")
	}
	if orderId := strings.TrimSpace(f.OrderId); orderId != "" {
		m = m.Where(string(entity.GameWinLogOrderId)+" LIKE ?", "%"+orderId+"%")
	}
	if platformType := strings.TrimSpace(f.PlatformType); platformType != "" {
		m = m.Where(string(entity.GameWinLogPlatformType)+" LIKE ?", "%"+platformType+"%")
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
