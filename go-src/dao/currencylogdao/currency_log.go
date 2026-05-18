package currencylogdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity"
)

// GetByUserId 按用户ID分页查询货币流水(按ID倒序)
func GetByUserId(userId uint64, pageIndex, pageSize int) (int, []*entity.CurrencyLog) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	ctx := gctx.New()
	list := make([]*entity.CurrencyLog, 0)
	m := g.Model(string(entity.TbCurrencyLog)).Ctx(ctx).
		Where(string(entity.CurrencyLogUserId)+" = ?", userId)
	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	_ = m.Clone().Order("id desc").
		Limit(pageSize).Offset((pageIndex - 1) * pageSize).
		Scan(&list)
	return total, list
}
