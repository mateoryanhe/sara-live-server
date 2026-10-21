package gift

import liveentity "xr-game-server/entity/live"

type luckyGiftMultiplierDefaultTier struct {
	Multiplier uint32
	Weight     uint32
	Type       string
}

var defaultLuckyGiftMultiplierTierDefs = []luckyGiftMultiplierDefaultTier{
	{Multiplier: 5, Weight: 50, Type: liveentity.LuckyGiftMultiplierTypeLow},
	{Multiplier: 10, Weight: 30, Type: liveentity.LuckyGiftMultiplierTypeLow},
	{Multiplier: 88, Weight: 15, Type: liveentity.LuckyGiftMultiplierTypeMedium},
	{Multiplier: 200, Weight: 4, Type: liveentity.LuckyGiftMultiplierTypeHigh},
	{Multiplier: 500, Weight: 1, Type: liveentity.LuckyGiftMultiplierTypeHigh},
}

func defaultLuckyGiftMultiplierTierRows() []*liveentity.LuckyGiftMultiplierTier {
	rows := make([]*liveentity.LuckyGiftMultiplierTier, 0, len(defaultLuckyGiftMultiplierTierDefs))
	for i, def := range defaultLuckyGiftMultiplierTierDefs {
		rows = append(rows, &liveentity.LuckyGiftMultiplierTier{
			Multiplier: def.Multiplier,
			Weight:     def.Weight,
			Type:       def.Type,
			Sort:       i,
		})
	}
	return rows
}
