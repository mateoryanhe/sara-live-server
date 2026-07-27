package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/accountcfgdto"
	"xr-game-server/module/accountcfg"
)

const AccountCfgCMSUrl = "/accountCfg"

type AccountCfgCMSController struct{}

func initAccountCfgCMSController() {
	httpserver.RegCMS(AccountCfgCMSUrl, &AccountCfgCMSController{})
}

func (c *AccountCfgCMSController) GetAccountCfg(ctx context.Context, req *accountcfgdto.GetAccountCfgReq) (*accountcfgdto.GetAccountCfgRes, error) {
	return accountcfg.GetAccountCfg(ctx, req)
}

func (c *AccountCfgCMSController) SaveAccountCfg(ctx context.Context, req *accountcfgdto.SaveAccountCfgReq) (*accountcfgdto.SaveAccountCfgRes, error) {
	return accountcfg.SaveAccountCfg(ctx, req)
}
