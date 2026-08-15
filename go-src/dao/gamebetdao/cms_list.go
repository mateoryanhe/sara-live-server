package gamebetdao

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/game"
)

// CMSListFilter CMS 下注记录查询条件
type CMSListFilter struct {
	UserId       uint64
	GameCode     string
	OrderId      string
	PlatformType string
	PageIndex    int
	PageSize     int
}

// CMSList CMS 分页查询游戏下注记录(按 ID 倒序)
func CMSList(f *CMSListFilter) (int, []*entity.GameBetLog) {
	list := make([]*entity.GameBetLog, 0)
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
	m := g.Model(string(entity.TbGameBetLog)).Ctx(ctx).
		Where(string(entity.GameBetLogAmount)+" > ?", 0)
	if f.UserId > 0 {
		m = m.Where(string(entity.GameBetLogUserId)+" = ?", f.UserId)
	}
	if gameCode := strings.TrimSpace(f.GameCode); gameCode != "" {
		m = m.Where(string(entity.GameBetLogGameCode)+" LIKE ?", "%"+gameCode+"%")
	}
	if orderId := strings.TrimSpace(f.OrderId); orderId != "" {
		m = m.Where(string(entity.GameBetLogOrderId)+" LIKE ?", "%"+orderId+"%")
	}
	if platformType := strings.TrimSpace(f.PlatformType); platformType != "" {
		m = m.Where(string(entity.GameBetLogPlatformType)+" LIKE ?", "%"+platformType+"%")
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
