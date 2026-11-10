package luckygiftpoolcfg

import liveentity "xr-game-server/entity/live"

type globalTierDef struct {
	Multiplier uint32
	Weight     uint32
	Type       string
}

// 与 CMS / 单礼物默认 5 挡保持一致
var defaultGlobalTierDefs = []globalTierDef{
	{Multiplier: 5, Weight: 50, Type: liveentity.LuckyGiftMultiplierTypeLow},
	{Multiplier: 10, Weight: 30, Type: liveentity.LuckyGiftMultiplierTypeLow},
	{Multiplier: 88, Weight: 15, Type: liveentity.LuckyGiftMultiplierTypeMedium},
	{Multiplier: 200, Weight: 4, Type: liveentity.LuckyGiftMultiplierTypeHigh},
	{Multiplier: 500, Weight: 1, Type: liveentity.LuckyGiftMultiplierTypeHigh},
}

const requiredGlobalTierCount = 5

func defaultGlobalTierRows() []*liveentity.LuckyGiftMultiplierTier {
	rows := make([]*liveentity.LuckyGiftMultiplierTier, 0, len(defaultGlobalTierDefs))
	for i, def := range defaultGlobalTierDefs {
		rows = append(rows, &liveentity.LuckyGiftMultiplierTier{
			GiftID:     liveentity.GlobalLuckyGiftMultiplierGiftID,
			Multiplier: def.Multiplier,
			Weight:     def.Weight,
			Type:       def.Type,
			Sort:       i,
		})
	}
	return rows
}
