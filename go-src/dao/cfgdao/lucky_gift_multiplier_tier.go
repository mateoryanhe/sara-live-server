package cfgdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/live"
)

func LoadLuckyGiftMultiplierTiersByGiftID(giftID uint64) []*entity.LuckyGiftMultiplierTier {
	var rows []*entity.LuckyGiftMultiplierTier
	if err := g.DB().Model(string(entity.TbLuckyGiftMultiplierTier)).
		Where("gift_id = ?", giftID).
		Order("sort asc, id asc").
		Scan(&rows); err != nil {
		return nil
	}
	return rows
}

// LoadGlobalLuckyGiftMultiplierTiers 全局默认倍率挡位(gift_id=0)
func LoadGlobalLuckyGiftMultiplierTiers() []*entity.LuckyGiftMultiplierTier {
	return LoadLuckyGiftMultiplierTiersByGiftID(entity.GlobalLuckyGiftMultiplierGiftID)
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

// ReplaceLuckyGiftMultiplierTiersByGiftID 同步删除旧挡位后，用 syndb New 全量写入
func ReplaceLuckyGiftMultiplierTiersByGiftID(giftID uint64, rows []*entity.LuckyGiftMultiplierTier) error {
	if err := DeleteLuckyGiftMultiplierTiersByGiftID(giftID); err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}
	for i, row := range rows {
		if row == nil {
			continue
		}
		_ = entity.NewLuckyGiftMultiplierTier(giftID, row.Multiplier, row.Weight, row.Type, i)
	}
	return nil
}

// ReplaceGlobalLuckyGiftMultiplierTiers 覆盖全局默认倍率挡位(gift_id=0)
func ReplaceGlobalLuckyGiftMultiplierTiers(rows []*entity.LuckyGiftMultiplierTier) error {
	return ReplaceLuckyGiftMultiplierTiersByGiftID(entity.GlobalLuckyGiftMultiplierGiftID, rows)
}

func DeleteLuckyGiftMultiplierTiersByGiftID(giftID uint64) error {
	_, err := g.DB().Model(string(entity.TbLuckyGiftMultiplierTier)).Where("gift_id = ?", giftID).Delete()
	return err
}
