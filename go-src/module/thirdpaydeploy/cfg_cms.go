package thirdpaydeploy

import (
	"context"
	"strings"
	"time"

	"xr-game-server/core/cfg"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/thirdpaydeploydto"
	rechargeentity "xr-game-server/entity/recharge"
	"xr-game-server/errercode"
	"xr-game-server/module/domainsite"
)

func thirdPayVisibilityFromHaiPay() (androidVisible, h5Visible bool) {
	androidVisible, h5Visible = true, true
	if cfg := cfgdao.GetHaiPayCfgCached(); cfg != nil {
		androidVisible = cfg.TVisable
		h5Visible = cfg.TH5Visable
	}
	return androidVisible, h5Visible
}

func saveThirdPayVisibility(androidVisible, h5Visible bool) error {
	existing := cfgdao.GetHaiPayCfgCached()
	now := time.Now()
	row := &rechargeentity.HaiPayCfg{
		TVisable:     androidVisible,
		TH5Visable:   h5Visible,
		UseProdAppID: true,
		Subject:      "Recharge",
	}
	if existing != nil {
		row.ID = existing.ID
		row.CreatedAt = existing.CreatedAt
		row.ApiHost = existing.ApiHost
		row.MerchantSecretKey = existing.MerchantSecretKey
		row.MerchantPrivateKey = existing.MerchantPrivateKey
		row.CallbackBaseUrl = existing.CallbackBaseUrl
		row.ReturnUrl = existing.ReturnUrl
		row.FailReturnUrl = existing.FailReturnUrl
		row.CancelUrl = existing.CancelUrl
		row.PaymentMethods = existing.PaymentMethods
		row.Subject = existing.Subject
		if row.Subject == "" {
			row.Subject = "Recharge"
		}
		row.UseProdAppID = existing.UseProdAppID
	}
	row.UpdatedAt = now
	if row.CreatedAt.IsZero() {
		row.CreatedAt = now
	}
	if err := cfgdao.SaveHaiPayCfg(row); err != nil {
		return err
	}
	cfgdao.ReloadHaiPayCfgCache()
	return nil
}

func GetThirdPayDeployInfo(ctx context.Context, _ *thirdpaydeploydto.GetThirdPayDeployInfoReq) (*thirdpaydeploydto.GetThirdPayDeployInfoRes, error) {
	deployPath, err := getDeployDir()
	if err != nil {
		return nil, err
	}
	androidVisible, h5Visible := thirdPayVisibilityFromHaiPay()
	return &thirdpaydeploydto.GetThirdPayDeployInfoRes{
		Info: &thirdpaydeploydto.ThirdPayDeployInfoItem{
			Domain:         cfg.GetStaticSiteDomains(thirdpaydeploydto.ThirdPayStaticPrefix),
			UrlPrefix:      thirdpaydeploydto.ThirdPayStaticPrefix,
			DeployPath:     deployPath,
			AcceptExt:      ".zip",
			LastUploadAt:   domainsite.GetLastUploadAt(ctx, thirdpaydeploydto.ThirdPaySiteKey),
			AndroidVisible: androidVisible,
			H5Visible:      h5Visible,
		},
	}, nil
}

func SaveThirdPayDeployCfg(ctx context.Context, req *thirdpaydeploydto.SaveThirdPayDeployCfgReq) (*thirdpaydeploydto.SaveThirdPayDeployCfgRes, error) {
	domain := strings.TrimSpace(req.Domain)
	deployPath := strings.TrimSpace(req.DeployPath)
	if domain == "" || deployPath == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if _, err := domainsite.SaveStaticSiteMapping(
		ctx,
		thirdpaydeploydto.ThirdPaySiteKey,
		thirdpaydeploydto.ThirdPayStaticPrefix,
		domain,
		deployPath,
	); err != nil {
		return nil, err
	}
	androidVisible, h5Visible := thirdPayVisibilityFromHaiPay()
	if req.AndroidVisible != nil {
		androidVisible = *req.AndroidVisible
	}
	if req.H5Visible != nil {
		h5Visible = *req.H5Visible
	}
	if err := saveThirdPayVisibility(androidVisible, h5Visible); err != nil {
		return nil, err
	}
	return &thirdpaydeploydto.SaveThirdPayDeployCfgRes{Success: true}, nil
}
