package cfgdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/live"
)

type LuckyGiftPoolPendingBalanceListFilter struct {
	Status    *uint8
	PageIndex int
	PageSize  int
}

// ListLuckyGiftPoolPendingBalances CMS 分页列表(新→旧)
func ListLuckyGiftPoolPendingBalances(f *LuckyGiftPoolPendingBalanceListFilter) (int, []*entity.LuckyGiftPoolPendingBalance) {
	if f == nil {
		f = &LuckyGiftPoolPendingBalanceListFilter{}
	}
	pageIndex := f.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	pageSize := f.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 200 {
		pageSize = 200
	}
	m := g.DB().Model(string(entity.TbLuckyGiftPoolPendingBalance))
	if f.Status != nil {
		m = m.Where("status = ?", *f.Status)
	}
	total, err := m.Count()
	if err != nil {
		return 0, nil
	}
	var rows []*entity.LuckyGiftPoolPendingBalance
	if err := m.Order("id desc").
		Limit(pageSize).
		Offset((pageIndex - 1) * pageSize).
		Scan(&rows); err != nil {
		return int(total), nil
	}
	return int(total), rows
}

// LoadPendingLuckyGiftPoolBalances 未处理的临时入账(按创建时间升序)
func LoadPendingLuckyGiftPoolBalances() []*entity.LuckyGiftPoolPendingBalance {
	var rows []*entity.LuckyGiftPoolPendingBalance
	if err := g.DB().Model(string(entity.TbLuckyGiftPoolPendingBalance)).
		Where("status = ?", entity.LuckyGiftPoolPendingStatusPending).
		Order("created_at asc, id asc").
		Scan(&rows); err != nil {
		return nil
	}
	return rows
}

// ClaimLuckyGiftPoolPendingBalance 乐观抢占: pending → processing, 防并发双加
func ClaimLuckyGiftPoolPendingBalance(id uint64) bool {
	if id == 0 {
		return false
	}
	r, err := g.DB().Model(string(entity.TbLuckyGiftPoolPendingBalance)).
		Where("id = ? AND status = ?", id, entity.LuckyGiftPoolPendingStatusPending).
		Data(g.Map{
			"status":     entity.LuckyGiftPoolPendingStatusProcessing,
			"updated_at": time.Now(),
		}).
		Update()
	if err != nil {
		return false
	}
	n, _ := r.RowsAffected()
	return n > 0
}

// MarkLuckyGiftPoolPendingBalanceProcessed 标记已并入奖池
func MarkLuckyGiftPoolPendingBalanceProcessed(id uint64) error {
	if id == 0 {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbLuckyGiftPoolPendingBalance)).
		Where("id = ?", id).
		Data(g.Map{
			"status":     entity.LuckyGiftPoolPendingStatusProcessed,
			"updated_at": time.Now(),
		}).
		Update()
	return err
}

// ResetStuckLuckyGiftPoolPendingProcessing 启动时把卡在 processing 的行复位为 pending
func ResetStuckLuckyGiftPoolPendingProcessing() {
	_, _ = g.DB().Model(string(entity.TbLuckyGiftPoolPendingBalance)).
		Where("status = ?", entity.LuckyGiftPoolPendingStatusProcessing).
		Data(g.Map{
			"status":     entity.LuckyGiftPoolPendingStatusPending,
			"updated_at": time.Now(),
		}).
		Update()
}
