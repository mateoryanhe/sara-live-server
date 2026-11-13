package luckygiftpoolcfg

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/dto/giftdto"
	"xr-game-server/dto/luckygiftpoolcfgdto"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
)

func GetLuckyGiftPoolCfg(_ context.Context, _ *luckygiftpoolcfgdto.GetLuckyGiftPoolCfgReq) (*luckygiftpoolcfgdto.GetLuckyGiftPoolCfgRes, error) {
	id, poolBalance, contributePercent, winProbability, totalConsume, totalPayout, tiers, createdAt, _ := snapshotForCMS()
	item := &luckygiftpoolcfgdto.LuckyGiftPoolCfgItem{
		PoolBalance:        poolBalance,
		ContributePercent:  contributePercent,
		WinProbability:     winProbability,
		TotalConsumeAmount: totalConsume,
		TotalPayoutAmount:  totalPayout,
		Tiers:              toTierItems(tiers),
		CreatedAt:          formatCfgTime(createdAt),
		UpdatedAt:          formatCfgTime(time.Now()),
	}
	if id != 0 {
		item.ID = strconv.FormatUint(id, 10)
	}
	if len(item.Tiers) == 0 {
		item.Tiers = toTierItems(defaultGlobalTierRows())
	}
	return &luckygiftpoolcfgdto.GetLuckyGiftPoolCfgRes{Cfg: item}, nil
}

func SaveLuckyGiftPoolCfg(_ context.Context, req *luckygiftpoolcfgdto.SaveLuckyGiftPoolCfgReq) (*luckygiftpoolcfgdto.SaveLuckyGiftPoolCfgRes, error) {
	if req.PoolBalance < 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if req.ContributePercent < 0 || req.ContributePercent > 100 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if req.WinProbability < 0 || req.WinProbability > 100 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if len(req.Tiers) != requiredGlobalTierCount {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	tierInputs := make([]*tierInput, 0, requiredGlobalTierCount)
	for _, item := range req.Tiers {
		if item == nil {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		tierInputs = append(tierInputs, &tierInput{
			Multiplier: item.Multiplier,
			Weight:     item.Weight,
			Type:       item.Type,
		})
	}

	var reqCfgID uint64
	if req.ID != "" {
		parsed, perr := strconv.ParseUint(req.ID, 10, 64)
		if perr != nil {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		reqCfgID = parsed
	}
	cfgID, err := enqueueApplyConfig(&applyConfigPayload{
		cfgID:             reqCfgID,
		poolBalance:       req.PoolBalance,
		contributePercent: req.ContributePercent,
		winProbability:    req.WinProbability,
		tiers:             tierInputs,
	})
	if err != nil {
		return nil, err
	}
	return &luckygiftpoolcfgdto.SaveLuckyGiftPoolCfgRes{
		Success: true,
		ID:      strconv.FormatUint(cfgID, 10),
	}, nil
}

func toTierItems(list []*liveentity.LuckyGiftMultiplierTier) []*giftdto.LuckyGiftMultiplierTierItem {
	if len(list) == 0 {
		return []*giftdto.LuckyGiftMultiplierTierItem{}
	}
	out := make([]*giftdto.LuckyGiftMultiplierTierItem, 0, len(list))
	for _, row := range list {
		if row == nil {
			continue
		}
		id := ""
		if row.ID != 0 {
			id = strconv.FormatUint(row.ID, 10)
		}
		out = append(out, &giftdto.LuckyGiftMultiplierTierItem{
			ID:         id,
			Multiplier: row.Multiplier,
			Weight:     row.Weight,
			Type:       row.Type,
		})
	}
	return out
}

func formatCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
