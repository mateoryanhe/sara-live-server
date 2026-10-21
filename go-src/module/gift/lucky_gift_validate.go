package gift

import (
	"xr-game-server/dto/giftdto"
	liveentity "xr-game-server/entity/live"
)

func applyLuckyGiftFields(g *liveentity.LiveGift, categoryType uint8, rtp float64) {
	if g == nil {
		return
	}
	g.CategoryType = categoryType
	if categoryType == liveentity.LiveGiftCategoryTypeLucky {
		g.LuckyGiftRtp = NormalizeLuckyGiftRtp(rtp)
		return
	}
	g.LuckyGiftRtp = 0
}

func validateLuckyGiftConfig(categoryType uint8, tiers []*giftdto.LuckyGiftMultiplierTierItem) error {
	if categoryType != liveentity.LiveGiftCategoryTypeLucky {
		return nil
	}
	if _, err := luckyGiftMultiplierRowsFromDTO(tiers); err != nil {
		return err
	}
	return nil
}
