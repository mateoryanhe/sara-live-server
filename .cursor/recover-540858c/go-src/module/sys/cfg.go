package sys

import (
	"context"
	"time"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/sysdto"
	"xr-game-server/module/apppkg"
	"xr-game-server/module/livecfg"
	"xr-game-server/module/privacypolicy"
	"xr-game-server/module/upload"
	"xr-game-server/module/wallet"
)

func resolveAboutSiteUrl() string {
	if url := privacypolicy.GetAboutSiteUrl(); url != "" {
		return url
	}
	return upload.GetAboutSiteUrl()
}

func resolveSafetyCenterUrl() string {
	if url := privacypolicy.GetSafetyCenterUrl(); url != "" {
		return url
	}
	return upload.GetSafetyCenterUrl()
}

func GetSysCfg(ctx context.Context, req *sysdto.SysCfgReq) (*sysdto.SysCfgResp, error) {
	_ = req
	packageName := httpserver.GetPackageNameFromContext(ctx)
	globalPrivacy := privacypolicy.GetPrivacyPolicyUrl()
	globalTerms := privacypolicy.GetTermsOfServiceUrl()
	exchangeCfg := wallet.GetExchangeCfgSnapshot()
	return &sysdto.SysCfgResp{
		SysTime:                     time.Now().UnixMilli(),
		PaidDanmakuPrice:            livecfg.GetPaidDanmakuPrice(),
		PrivateRoomFreeWatchSeconds: livecfg.GetPrivateRoomFreeWatchSeconds(),
		PrivacyPolicyUrl:            apppkg.ResolvePrivacyPolicyUrl(packageName, globalPrivacy),
		TermsOfServiceUrl:           apppkg.ResolveTermsOfServiceUrl(packageName, globalTerms),
		CreatorTermsUrl:             privacypolicy.GetCreatorTermsUrl(),
		RoomOwnerTermsUrl:           privacypolicy.GetRoomOwnerTermsUrl(),
		VipDescUrl:                  privacypolicy.GetVipDescUrl(),
		AppImageMaxSize:             upload.GetAppImageMaxSize(),
		GoldToDiamondRate:           exchangeCfg.GoldToDiamondRate,
		ExchangeFeePercent:          exchangeCfg.ExchangeFeePercent,
		UsdToGoldRate:               exchangeCfg.UsdToGoldRate,
		AboutSiteUrl:                resolveAboutSiteUrl(),
		SafetyCenterUrl:             resolveSafetyCenterUrl(),
	}, nil
}
