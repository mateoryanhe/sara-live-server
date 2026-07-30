package customerservice

import (
	"context"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/customerservicedto"
)

func GetAppCustomerServiceCfg(_ context.Context, _ *customerservicedto.AppCustomerServiceCfgReq) (*customerservicedto.AppCustomerServiceCfgRes, error) {
	cfg := cfgdao.GetCustomerServiceCfgCached()
	if cfg == nil {
		return &customerservicedto.AppCustomerServiceCfgRes{}, nil
	}
	return &customerservicedto.AppCustomerServiceCfgRes{
		TelegramUrl: cfg.TelegramUrl,
		FacebookUrl: cfg.FacebookUrl,
		WhatsappUrl: cfg.WhatsappUrl,
	}, nil
}
