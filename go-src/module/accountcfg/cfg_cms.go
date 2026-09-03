package accountcfg

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/accountcfgdto"
	"xr-game-server/entity/cms"
	"xr-game-server/errercode"
)

func GetAccountCfg(_ context.Context, _ *accountcfgdto.GetAccountCfgReq) (*accountcfgdto.GetAccountCfgRes, error) {
	cfg := cfgdao.LoadAccountCfg()
	if cfg == nil {
		return &accountcfgdto.GetAccountCfgRes{Cfg: &accountcfgdto.AccountCfgItem{
			CancelAccountByCodeEnabled: false,
			BlockSimulatorLogin:        false,
			EnvType:                    entity.AccountEnvTypeProd,
		}}, nil
	}
	return &accountcfgdto.GetAccountCfgRes{Cfg: toCfgItem(cfg)}, nil
}

func SaveAccountCfg(_ context.Context, req *accountcfgdto.SaveAccountCfgReq) (*accountcfgdto.SaveAccountCfgRes, error) {
	if req.EnvType != entity.AccountEnvTypeProd &&
		req.EnvType != entity.AccountEnvTypeReview &&
		req.EnvType != entity.AccountEnvTypeTest {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	existing := cfgdao.LoadAccountCfg()
	row := &entity.AccountCfg{
		CancelAccountByCodeEnabled: req.CancelAccountByCodeEnabled,
		BlockSimulatorLogin:        req.BlockSimulatorLogin,
		EnvType:                    req.EnvType,
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
	if err := cfgdao.SaveAccountCfg(row); err != nil {
		return nil, err
	}
	reloadCfgMemory()
	return &accountcfgdto.SaveAccountCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toCfgItem(cfg *entity.AccountCfg) *accountcfgdto.AccountCfgItem {
	if cfg == nil {
		return nil
	}
	return &accountcfgdto.AccountCfgItem{
		ID:                         strconv.FormatUint(cfg.ID, 10),
		CancelAccountByCodeEnabled: cfg.CancelAccountByCodeEnabled,
		BlockSimulatorLogin:        cfg.BlockSimulatorLogin,
		EnvType:                    cfg.EnvType,
		CreatedAt:                  formatTime(cfg.CreatedAt),
		UpdatedAt:                  formatTime(cfg.UpdatedAt),
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
