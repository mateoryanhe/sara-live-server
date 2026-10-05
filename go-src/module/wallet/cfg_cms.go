package wallet

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/walletdto"
	"xr-game-server/entity/user"
	"xr-game-server/errercode"
)

func GetWalletExchangeCfg(_ context.Context, _ *walletdto.GetWalletExchangeCfgReq) (*walletdto.GetWalletExchangeCfgRes, error) {
	cfg := cfgdao.GetWalletExchangeCfgCached()
	if cfg == nil {
		return &walletdto.GetWalletExchangeCfgRes{
			Cfg: &walletdto.WalletExchangeCfgItem{
				GoldToDiamondRate:  DefaultGoldToDiamondRate,
				ExchangeFeePercent: DefaultExchangeFeePercent,
				UsdToGoldRate:      DefaultUsdToGoldRate,
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
	if req.UsdToGoldRate <= 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	row, matched, err := cfgdao.UpdateWalletExchangeCfg(req.ID, func(row *entity.WalletExchangeCfg) {
		row.GoldToDiamondRate = req.GoldToDiamondRate
		row.ExchangeFeePercent = req.ExchangeFeePercent
		row.UsdToGoldRate = req.UsdToGoldRate
	})
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return &walletdto.SaveWalletExchangeCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toWalletExchangeCfgItem(cfg *entity.WalletExchangeCfg) *walletdto.WalletExchangeCfgItem {
	if cfg == nil {
		return nil
	}
	usdToGold := cfg.UsdToGoldRate
	if usdToGold <= 0 {
		usdToGold = DefaultUsdToGoldRate
	}
	return &walletdto.WalletExchangeCfgItem{
		ID:                 strconv.FormatUint(cfg.ID, 10),
		GoldToDiamondRate:  cfg.GoldToDiamondRate,
		ExchangeFeePercent: cfg.ExchangeFeePercent,
		UsdToGoldRate:      usdToGold,
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
