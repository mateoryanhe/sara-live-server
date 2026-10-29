package sys

import (
	"context"
	"time"

	corecfg "xr-game-server/core/cfg"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/sysdto"
	"xr-game-server/module/accountcfg"
	"xr-game-server/module/apppkg"
	"xr-game-server/module/livecfg"
	"xr-game-server/module/privacypolicy"
	"xr-game-server/module/upload"
	"xr-game-server/dto/userwithdrawdeploydto"
	"xr-game-server/module/wallet"
)

func resolveAboutSiteUrl(useSofie bool) string {
	if useSofie {
		if url := privacypolicy.GetSofieAboutSiteUrl(); url != "" {
			return url
		}
		// Sofie 未填时回落 Sara
	}
	if url := privacypolicy.GetAboutSiteUrl(); url != "" {
		return url
	}
	return upload.GetAboutSiteUrl()
}

func resolveSafetyCenterUrl(useSofie bool) string {
	if useSofie {
		if url := privacypolicy.GetSofieSafetyCenterUrl(); url != "" {
			return url
		}
		// Sofie 未填时回落 Sara
	}
	if url := privacypolicy.GetSafetyCenterUrl(); url != "" {
		return url
	}
	return upload.GetSafetyCenterUrl()
}

func GetSysCfg(ctx context.Context, req *sysdto.SysCfgReq) (*sysdto.SysCfgResp, error) {
	_ = req
	packageName := httpserver.GetPackageNameFromContext(ctx)
	sofie := httpserver.IsSofieAppRequest(ctx)
	useSofiePrivacy := privacypolicy.UseSofiePrivacyForPackage(sofie, packageName)
	var (
		privacyUrl, termsUrl, creatorTerms, roomOwnerTerms, vipDesc string
	)
	if useSofiePrivacy {
		privacyUrl = privacypolicy.GetSofiePrivacyPolicyUrl()
		termsUrl = privacypolicy.GetSofieTermsOfServiceUrl()
		creatorTerms = privacypolicy.GetSofieCreatorTermsUrl()
		roomOwnerTerms = privacypolicy.GetSofieRoomOwnerTermsUrl()
		vipDesc = privacypolicy.GetSofieVipDescUrl()
		// 单字段为空时回落 Sara 隐私政策
		if privacyUrl == "" {
			privacyUrl = privacypolicy.GetPrivacyPolicyUrl()
		}
		if termsUrl == "" {
			termsUrl = privacypolicy.GetTermsOfServiceUrl()
		}
		if creatorTerms == "" {
			creatorTerms = privacypolicy.GetCreatorTermsUrl()
		}
		if roomOwnerTerms == "" {
			roomOwnerTerms = privacypolicy.GetRoomOwnerTermsUrl()
		}
		if vipDesc == "" {
			vipDesc = privacypolicy.GetVipDescUrl()
		}
	} else {
		privacyUrl = apppkg.ResolvePrivacyPolicyUrl(packageName, privacypolicy.GetPrivacyPolicyUrl())
		termsUrl = apppkg.ResolveTermsOfServiceUrl(packageName, privacypolicy.GetTermsOfServiceUrl())
		creatorTerms = privacypolicy.GetCreatorTermsUrl()
		roomOwnerTerms = privacypolicy.GetRoomOwnerTermsUrl()
		vipDesc = privacypolicy.GetVipDescUrl()
	}
	exchangeCfg := wallet.GetExchangeCfgSnapshot()
	// 无 HaiPay 配置时默认显示;有配置则读库
	tVisable, tH5Visable := true, true
	if haiPayCfg := cfgdao.GetHaiPayCfgCached(); haiPayCfg != nil {
		tVisable = haiPayCfg.TVisable
		tH5Visable = haiPayCfg.TH5Visable
	}
	resp := &sysdto.SysCfgResp{
		SysTime:                    time.Now().UnixMilli(),
		T:                          corecfg.GetThirdPayDomain(),
		TVisable:                   tVisable,
		TH5Visable:                 tH5Visable,
		W:                          corecfg.GetStaticSiteDomain(userwithdrawdeploydto.UserWithdrawStaticPrefix),
		WVisable:                   accountcfg.IsUserWithdrawEnabled(),
		WAndroidVisable:            accountcfg.IsUserWithdrawAndroidVisible(),
		PaidDanmakuPrice:           livecfg.GetPaidDanmakuPrice(),
		AudienceListRefreshSeconds: livecfg.GetAudienceListRefreshSeconds(),
		OneToOneDailyFreeSeconds:   livecfg.GetOneToOneDailyFreeSeconds(),
		ShowcaseLiveRoomEnabled:    livecfg.IsShowcaseLiveRoomEnabled(),
		GameLiveRoomEnabled:        livecfg.IsGameLiveRoomEnabled(),
		VoiceChatLiveRoomEnabled:   livecfg.IsVoiceChatLiveRoomEnabled(),
		OneToOneLiveRoomEnabled:    livecfg.IsOneToOneLiveRoomEnabled(),
		PrivacyPolicyUrl:           privacyUrl,
		TermsOfServiceUrl:          termsUrl,
		CreatorTermsUrl:            creatorTerms,
		RoomOwnerTermsUrl:          roomOwnerTerms,
		VipDescUrl:                 vipDesc,
		AppImageMaxSize:            upload.GetAppImageMaxSize(),
		GoldToDiamondRate:          exchangeCfg.GoldToDiamondRate,
		ExchangeFeePercent:         exchangeCfg.ExchangeFeePercent,
		UsdToGoldRate:              exchangeCfg.UsdToGoldRate,
		AboutSiteUrl:               resolveAboutSiteUrl(useSofiePrivacy),
		SafetyCenterUrl:            resolveSafetyCenterUrl(useSofiePrivacy),
		EnvType:                    accountcfg.GetEnvType(),
	}
	if pkg := apppkg.GetAppPkgFromMemoryByPackageName(packageName); pkg != nil {
		resp.AttributionEnabled = pkg.AttributionEnabled
		resp.AttributionProvider = pkg.AttributionProvider
		resp.AppsFlyerDevKey = pkg.AppsFlyerDevKey
		resp.AppsFlyerAppId = pkg.AppsFlyerAppId
	}
	return resp, nil
}
