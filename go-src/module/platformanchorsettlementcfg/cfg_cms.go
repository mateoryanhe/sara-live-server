package platformanchorsettlementcfg

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/platformanchorsettlementcfgdto"
	"xr-game-server/entity/user"
	"xr-game-server/errercode"
)

func GetPlatformAnchorSettlementCfg(_ context.Context, _ *platformanchorsettlementcfgdto.GetPlatformAnchorSettlementCfgReq) (*platformanchorsettlementcfgdto.GetPlatformAnchorSettlementCfgRes, error) {
	cfg := cfgdao.GetWalletExchangeCfgCached()
	if cfg == nil {
		return &platformanchorsettlementcfgdto.GetPlatformAnchorSettlementCfgRes{
			Cfg: &platformanchorsettlementcfgdto.PlatformAnchorSettlementCfgItem{
				MinimumSettlementUsd: cfgdao.DefaultPlatformAnchorMinimumSettlementUsd,
			},
		}, nil
	}
	return &platformanchorsettlementcfgdto.GetPlatformAnchorSettlementCfgRes{Cfg: toPlatformAnchorSettlementCfgItem(cfg)}, nil
}

func SavePlatformAnchorSettlementCfg(_ context.Context, req *platformanchorsettlementcfgdto.SavePlatformAnchorSettlementCfgReq) (*platformanchorsettlementcfgdto.SavePlatformAnchorSettlementCfgRes, error) {
	if req.MinimumSettlementUsd <= 0 || req.MinimumSettlementUsd > cfgdao.MaxPlatformAnchorMinimumSettlementUsd {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	row, matched, err := cfgdao.UpdateWalletExchangeCfg(req.ID, func(row *entity.WalletExchangeCfg) {
		row.PlatformAnchorMinimumSettlementUsd = req.MinimumSettlementUsd
	})
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return &platformanchorsettlementcfgdto.SavePlatformAnchorSettlementCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toPlatformAnchorSettlementCfgItem(cfg *entity.WalletExchangeCfg) *platformanchorsettlementcfgdto.PlatformAnchorSettlementCfgItem {
	if cfg == nil {
		return nil
	}
	return &platformanchorsettlementcfgdto.PlatformAnchorSettlementCfgItem{
		ID:                   strconv.FormatUint(cfg.ID, 10),
		MinimumSettlementUsd: cfgdao.PlatformAnchorMinimumSettlementUsd(),
		CreatedAt:            formatCfgTime(cfg.CreatedAt),
		UpdatedAt:            formatCfgTime(cfg.UpdatedAt),
	}
}

func formatCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
