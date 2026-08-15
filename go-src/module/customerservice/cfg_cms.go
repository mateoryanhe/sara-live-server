package customerservice

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/customerservicedto"
	"xr-game-server/entity/user"
	"xr-game-server/errercode"
)

func GetCustomerServiceCfg(_ context.Context, _ *customerservicedto.GetCustomerServiceCfgReq) (*customerservicedto.GetCustomerServiceCfgRes, error) {
	cfg := cfgdao.GetCustomerServiceCfgCached()
	if cfg == nil {
		return &customerservicedto.GetCustomerServiceCfgRes{Cfg: nil}, nil
	}
	return &customerservicedto.GetCustomerServiceCfgRes{Cfg: toCfgItem(cfg)}, nil
}

func SaveCustomerServiceCfg(_ context.Context, req *customerservicedto.SaveCustomerServiceCfgReq) (*customerservicedto.SaveCustomerServiceCfgRes, error) {
	telegramUrl := strings.TrimSpace(req.TelegramUrl)
	facebookUrl := strings.TrimSpace(req.FacebookUrl)
	whatsappUrl := strings.TrimSpace(req.WhatsappUrl)

	existing := cfgdao.GetCustomerServiceCfgCached()
	row := &entity.CustomerServiceCfg{
		TelegramUrl: telegramUrl,
		FacebookUrl: facebookUrl,
		WhatsappUrl: whatsappUrl,
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
	if err := cfgdao.SaveCustomerServiceCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadCustomerServiceCfgCache()
	return &customerservicedto.SaveCustomerServiceCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func toCfgItem(cfg *entity.CustomerServiceCfg) *customerservicedto.CustomerServiceCfgItem {
	if cfg == nil {
		return nil
	}
	return &customerservicedto.CustomerServiceCfgItem{
		ID:          strconv.FormatUint(cfg.ID, 10),
		TelegramUrl: cfg.TelegramUrl,
		FacebookUrl: cfg.FacebookUrl,
		WhatsappUrl: cfg.WhatsappUrl,
		CreatedAt:   formatTime(cfg.CreatedAt),
		UpdatedAt:   formatTime(cfg.UpdatedAt),
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
