package gift

import (
	"context"
	"strconv"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/giftdto"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
)

func GetGiftLuckyCfg(_ context.Context, req *giftdto.GetGiftLuckyCfgReq) (*giftdto.GetGiftLuckyCfgRes, error) {
	g := cfgdao.GetGiftById(req.GiftID)
	if g == nil {
		return nil, errercode.CreateCode(errercode.GiftNonExist)
	}
	rows := cfgdao.LoadLuckyGiftMultiplierTiersByGiftID(req.GiftID)
	return &giftdto.GetGiftLuckyCfgRes{
		LuckyGiftRtp: NormalizeLuckyGiftRtp(g.LuckyGiftRtp),
		Tiers:        toLuckyGiftMultiplierTierItems(rows),
	}, nil
}

func saveGiftLuckyConfig(giftID uint64, categoryType uint8, rtp float64, tiers []*giftdto.LuckyGiftMultiplierTierItem) error {
	if categoryType != liveentity.LiveGiftCategoryTypeLucky {
		if err := cfgdao.DeleteLuckyGiftMultiplierTiersByGiftID(giftID); err != nil {
			return err
		}
		return nil
	}
	rows, err := luckyGiftMultiplierRowsFromDTO(tiers)
	if err != nil {
		return err
	}
	return cfgdao.ReplaceLuckyGiftMultiplierTiersByGiftID(giftID, rows)
}

func luckyGiftMultiplierRowsFromDTO(list []*giftdto.LuckyGiftMultiplierTierItem) ([]*liveentity.LuckyGiftMultiplierTier, error) {
	if len(list) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	rows := make([]*liveentity.LuckyGiftMultiplierTier, 0, len(list))
	for i, item := range list {
		if item == nil {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		if item.Multiplier == 0 {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		if item.Weight == 0 {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		if !liveentity.IsValidLuckyGiftMultiplierType(item.Type) {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		rows = append(rows, &liveentity.LuckyGiftMultiplierTier{
			Multiplier: item.Multiplier,
			Weight:     item.Weight,
			Type:       item.Type,
			Sort:       i,
		})
	}
	return rows, nil
}

func toLuckyGiftMultiplierTierItems(list []*liveentity.LuckyGiftMultiplierTier) []*giftdto.LuckyGiftMultiplierTierItem {
	if len(list) == 0 {
		return []*giftdto.LuckyGiftMultiplierTierItem{}
	}
	out := make([]*giftdto.LuckyGiftMultiplierTierItem, 0, len(list))
	for _, row := range list {
		if row == nil {
			continue
		}
		out = append(out, &giftdto.LuckyGiftMultiplierTierItem{
			ID:         formatGiftUintID(row.ID),
			Multiplier: row.Multiplier,
			Weight:     row.Weight,
			Type:       row.Type,
		})
	}
	return out
}

func formatGiftUintID(id uint64) string {
	if id == 0 {
		return ""
	}
	return strconv.FormatUint(id, 10)
}
