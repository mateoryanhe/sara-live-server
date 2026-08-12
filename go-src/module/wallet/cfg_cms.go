package wallet

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/walletdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func GetWalletExchangeCfg(_ context.Context, _ *walletdto.GetWalletExchangeCfgReq) (*walletdto.GetWalletExchangeCfgRes, error) {
	cfg := cfgdao.GetWalletExchangeCfgCached()
	if cfg == nil {
		return &walletdto.GetWalletExchangeCfgRes{
			Cfg: &walletdto.WalletExchangeCfgItem{
				GoldToDiamondRate:  DefaultGoldToDiamondRate,
				ExchangeFeePercent: DefaultExchangeFeePercent,
			},
		}, nil
	}
	return &walletdto.GetWalletExchangeCfgRes{Cfg: toWalletExchangeCfgItem(cfg)}, nil
}

func SaveWalletExchangeCfg(_ context.Context, req *walletdto.SaveWalletExchangeCfgReq) (*walletdto.SaveWalletExchangeCfgRes, error) {
	if req.GoldToDiamondRate <= 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if req.ExchangeFeePercent < 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	existing := cfgdao.GetWalletExchangeCfgCached()
	row := &entity.WalletExchangeCfg{
		GoldToDiamondRate:  req.GoldToDiamondRate,
		ExchangeFeePercent: req.ExchangeFeePercent,
	}
	if req.ID > 0 {
		if existing == nil || existing.ID != req.ID {
			return nil, errercode.CreateCode(errercode.InvalidParam)
		}
		row.ID = req.ID
		row.CreatedAt = existing.CreatedAt
	} else if existing != nil {
		row.ID = existing.ID
		row.CreatedAt = existing.CreatedAt
	}
	row.UpdatedAt = time.Now()
	if row.CreatedAt.IsZero() {
		row.CreatedAt = row.UpdatedAt
	}
	if err := cfgdao.SaveWalletExchangeCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadWalletExchangeCfgCache()
	return &walletdto.SaveWalletExchangeCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toWalletExchangeCfgItem(cfg *entity.WalletExchangeCfg) *walletdto.WalletExchangeCfgItem {
	if cfg == nil {
		return nil
	}
	return &walletdto.WalletExchangeCfgItem{
		ID:                 strconv.FormatUint(cfg.ID, 10),
		GoldToDiamondRate:  cfg.GoldToDiamondRate,
		ExchangeFeePercent: cfg.ExchangeFeePercent,
		CreatedAt:          formatCfgTime(cfg.CreatedAt),
		UpdatedAt:          formatCfgTime(cfg.UpdatedAt),
	}
}

func formatCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
