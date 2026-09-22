package coinmerchantdeploy

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/core/cfg"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/coinmerchantdeploydto"
	"xr-game-server/entity/sys"
	"xr-game-server/errercode"
	"xr-game-server/module/domainsite"
)

func GetCoinMerchantDeployInfo(ctx context.Context, _ *coinmerchantdeploydto.GetCoinMerchantDeployInfoReq) (*coinmerchantdeploydto.GetCoinMerchantDeployInfoRes, error) {
	deployPath, err := getDeployDir()
	if err != nil {
		return nil, err
	}
	snap := getCfgCache()
	return &coinmerchantdeploydto.GetCoinMerchantDeployInfoRes{
		Info: &coinmerchantdeploydto.CoinMerchantDeployInfoItem{
			ID:           formatUintID(snap.ID),
			Domain:       cfg.GetStaticSiteDomains(coinmerchantdeploydto.CoinMerchantStaticPrefix),
			UrlPrefix:    coinmerchantdeploydto.CoinMerchantStaticPrefix,
			DeployPath:   deployPath,
			AcceptExt:    ".zip",
			DeploySecret: snap.DeploySecret,
			UpdatedAt:    snap.UpdatedAt,
			LastUploadAt: domainsite.GetLastUploadAt(ctx, coinmerchantdeploydto.CoinMerchantSiteKey),
		},
	}, nil
}

func SaveCoinMerchantDeployCfg(ctx context.Context, req *coinmerchantdeploydto.SaveCoinMerchantDeployCfgReq) (*coinmerchantdeploydto.SaveCoinMerchantDeployCfgRes, error) {
	secret := strings.TrimSpace(req.DeploySecret)
	domain := strings.TrimSpace(req.Domain)
	deployPath := strings.TrimSpace(req.DeployPath)
	if secret == "" || domain == "" || deployPath == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	existing := cfgdao.LoadCoinMerchantDeployCfg()
	row := &entity.CoinMerchantDeployCfg{
		DeploySecret: secret,
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
	if _, err := domainsite.SaveStaticSiteMapping(
		ctx,
		coinmerchantdeploydto.CoinMerchantSiteKey,
		coinmerchantdeploydto.CoinMerchantStaticPrefix,
		domain,
		deployPath,
	); err != nil {
		return nil, err
	}
	if err := cfgdao.SaveCoinMerchantDeployCfg(row); err != nil {
		return nil, err
	}
	ReloadCoinMerchantDeployCache()
	return &coinmerchantdeploydto.SaveCoinMerchantDeployCfgRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

func formatUintID(id uint64) string {
	if id == 0 {
		return "0"
	}
	return strconv.FormatUint(id, 10)
}
