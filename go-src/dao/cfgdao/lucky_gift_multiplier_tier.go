package cfgdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/snowflake"
	"xr-game-server/entity/live"
)

func LoadLuckyGiftMultiplierTiersByGiftID(giftID uint64) []*entity.LuckyGiftMultiplierTier {
	if giftID == 0 {
		return nil
	}
	var rows []*entity.LuckyGiftMultiplierTier
	if err := g.DB().Model(string(entity.TbLuckyGiftMultiplierTier)).
		Where("gift_id = ?", giftID).
		Order("sort asc, id asc").
		Scan(&rows); err != nil {
		return nil
	}
	return rows
}

// LoadLuckyGiftMultiplierTiersByGiftIDs 批量加载所选礼物的幸运返奖倍率挡位(数据同步用)
func LoadLuckyGiftMultiplierTiersByGiftIDs(giftIDs []uint64) []*entity.LuckyGiftMultiplierTier {
	ids := make([]uint64, 0, len(giftIDs))
	seen := make(map[uint64]struct{}, len(giftIDs))
	for _, id := range giftIDs {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return []*entity.LuckyGiftMultiplierTier{}
	}
	var rows []*entity.LuckyGiftMultiplierTier
	if err := g.DB().Model(string(entity.TbLuckyGiftMultiplierTier)).
		Where("gift_id IN (?)", ids).
		Order("gift_id asc, sort asc, id asc").
		Scan(&rows); err != nil {
		return []*entity.LuckyGiftMultiplierTier{}
	}
	if rows == nil {
		return []*entity.LuckyGiftMultiplierTier{}
	}
	return rows
}

func LoadAllLuckyGiftMultiplierTiers() []*entity.LuckyGiftMultiplierTier {
	var rows []*entity.LuckyGiftMultiplierTier
	if err := g.DB().Model(string(entity.TbLuckyGiftMultiplierTier)).
		Where("gift_id > ?", 0).
		Order("gift_id asc, sort asc, id asc").
		Scan(&rows); err != nil {
		return nil
	}
	return rows
}

// ReplaceLuckyGiftMultiplierTiersByGiftID 删除该礼物旧挡位后全量写入
func ReplaceLuckyGiftMultiplierTiersByGiftID(giftID uint64, rows []*entity.LuckyGiftMultiplierTier) error {
	if giftID == 0 {
		return nil
	}
	if _, err := g.DB().Model(string(entity.TbLuckyGiftMultiplierTier)).Where("gift_id = ?", giftID).Delete(); err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}
	now := time.Now()
	for _, row := range rows {
		if row == nil {
			continue
		}
		row.GiftID = giftID
		if row.ID == 0 {
			row.ID = snowflake.GetId()
		}
		if row.CreatedAt.IsZero() {
			row.CreatedAt = now
		}
		row.UpdatedAt = now
	}
	_, err := g.DB().Model(string(entity.TbLuckyGiftMultiplierTier)).Data(rows).Batch(len(rows)).Insert()
	return err
}

func DeleteLuckyGiftMultiplierTiersByGiftID(giftID uint64) error {
	if giftID == 0 {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbLuckyGiftMultiplierTier)).Where("gift_id = ?", giftID).Delete()
	return err
}
