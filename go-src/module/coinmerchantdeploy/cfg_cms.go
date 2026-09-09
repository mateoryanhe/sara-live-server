package coinmerchantdeploy

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/coinmerchantdeploydto"
	"xr-game-server/entity/sys"
	"xr-game-server/errercode"
)

func GetCoinMerchantDeployInfo(_ context.Context, _ *coinmerchantdeploydto.GetCoinMerchantDeployInfoReq) (*coinmerchantdeploydto.GetCoinMerchantDeployInfoRes, error) {
	deployPath, err := getDeployDir()
	if err != nil {
		return nil, err
	}
	snap := getCfgCache()
	return &coinmerchantdeploydto.GetCoinMerchantDeployInfoRes{
		Info: &coinmerchantdeploydto.CoinMerchantDeployInfoItem{
			ID:           formatUintID(snap.ID),
			UrlPrefix:    coinmerchantdeploydto.CoinMerchantStaticPrefix,
			DeployPath:   deployPath,
			AcceptExt:    ".zip",
			DeploySecret: snap.DeploySecret,
			UpdatedAt:    snap.UpdatedAt,
		},
	}, nil
}

func SaveCoinMerchantDeployCfg(_ context.Context, req *coinmerchantdeploydto.SaveCoinMerchantDeployCfgReq) (*coinmerchantdeploydto.SaveCoinMerchantDeployCfgRes, error) {
	secret := strings.TrimSpace(req.DeploySecret)
	if secret == "" {
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
