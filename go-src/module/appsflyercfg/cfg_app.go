package appsflyercfg

import (
	"context"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/appsflyercfgdto"
	sysentity "xr-game-server/entity/sys"
)

func GetAppsFlyerCfgForApp(ctx context.Context, _ *appsflyercfgdto.GetAppsFlyerCfgForAppReq) (*appsflyercfgdto.GetAppsFlyerCfgForAppRes, error) {
	cfg := cfgdao.GetAppsFlyerCfgCached()
	if cfg == nil {
		return &appsflyercfgdto.GetAppsFlyerCfgForAppRes{}, nil
	}
	platform := httpserver.ClientPlatformFromContext(ctx)
	item := resolveAppsFlyerCfgForPlatform(cfg, platform)
	if item == nil {
		return &appsflyercfgdto.GetAppsFlyerCfgForAppRes{}, nil
	}
	return &appsflyercfgdto.GetAppsFlyerCfgForAppRes{Cfg: item}, nil
}

func resolveAppsFlyerCfgForPlatform(cfg *sysentity.AppsFlyerCfg, platform string) *appsflyercfgdto.AppsFlyerCfgForApp {
	if cfg == nil {
		return nil
	}
	switch platform {
	case "android":
		devKey := strings.TrimSpace(cfg.AndroidDevKey)
		if devKey == "" {
			return nil
		}
		return &appsflyercfgdto.AppsFlyerCfgForApp{AppsFlyerDevKey: devKey}
	case "ios":
		devKey := strings.TrimSpace(cfg.IosDevKey)
		if devKey == "" {
			return nil
		}
		return &appsflyercfgdto.AppsFlyerCfgForApp{
			AppsFlyerDevKey: devKey,
			AppsFlyerAppId:  strings.TrimSpace(cfg.IosAppId),
		}
	default:
		return nil
	}
}