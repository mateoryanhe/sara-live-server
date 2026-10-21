package luckygiftsettlementcfg

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/luckygiftsettlementcfgdto"
	"xr-game-server/entity/user"
	"xr-game-server/errercode"
)

func GetLuckyGiftSettlementCfg(_ context.Context, _ *luckygiftsettlementcfgdto.GetLuckyGiftSettlementCfgReq) (*luckygiftsettlementcfgdto.GetLuckyGiftSettlementCfgRes, error) {
	cfg := cfgdao.GetWalletExchangeCfgCached()
	if cfg == nil {
		return &luckygiftsettlementcfgdto.GetLuckyGiftSettlementCfgRes{
			Cfg: &luckygiftsettlementcfgdto.LuckyGiftSettlementCfgItem{
				SettlementPercent: cfgdao.DefaultLuckyGiftSettlementPercent,
			},
		}, nil
	}
	return &luckygiftsettlementcfgdto.GetLuckyGiftSettlementCfgRes{Cfg: toLuckyGiftSettlementCfgItem(cfg)}, nil
}

func SaveLuckyGiftSettlementCfg(_ context.Context, req *luckygiftsettlementcfgdto.SaveLuckyGiftSettlementCfgReq) (*luckygiftsettlementcfgdto.SaveLuckyGiftSettlementCfgRes, error) {
	if req.SettlementPercent < 0 || req.SettlementPercent > cfgdao.MaxLuckyGiftSettlementPercent {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	row, matched, err := cfgdao.UpdateWalletExchangeCfg(req.ID, func(row *entity.WalletExchangeCfg) {
		row.LuckyGiftSettlementPercent = req.SettlementPercent
	})
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return &luckygiftsettlementcfgdto.SaveLuckyGiftSettlementCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toLuckyGiftSettlementCfgItem(cfg *entity.WalletExchangeCfg) *luckygiftsettlementcfgdto.LuckyGiftSettlementCfgItem {
	if cfg == nil {
		return nil
	}
	return &luckygiftsettlementcfgdto.LuckyGiftSettlementCfgItem{
		ID:                strconv.FormatUint(cfg.ID, 10),
		SettlementPercent: cfgdao.LuckyGiftSettlementPercent(),
		CreatedAt:         formatCfgTime(cfg.CreatedAt),
		UpdatedAt:         formatCfgTime(cfg.UpdatedAt),
	}
}

func formatCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
