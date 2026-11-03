package stat

import (
	"strings"

	"xr-game-server/gameevent"
	"xr-game-server/module/promotionattribution"
)

func resolveApkChannelCodeForAppRegister(val *gameevent.RegisterEventData) string {
	if val == nil {
		return ""
	}
	if apk := strings.TrimSpace(val.ApkPromotionChannel); apk != "" {
		return apk
	}
	promo := appRegisterPromotionURLChannelCode(val)
	if promo <= 0 {
		return ""
	}
	apk, ok := promotionattribution.ResolveH5AppPairByPromotionChannelCode(promo)
	if !ok {
		return ""
	}
	return apk
}

func appRegisterPromotionURLChannelCode(val *gameevent.RegisterEventData) int {
	if val == nil {
		return 0
	}
	if val.PairedPromotionChannelCode > 0 {
		return val.PairedPromotionChannelCode
	}
	if apk := strings.TrimSpace(val.ApkPromotionChannel); apk != "" {
		if c, ok := promotionattribution.ResolvePromotionChannelCodeByApkChannelCode(apk); ok {
			return c
		}
	}
	if val.PromotionChannelCode > 0 {
		if _, ok := promotionattribution.ResolveH5AppPairByPromotionChannelCode(val.PromotionChannelCode); ok {
			return val.PromotionChannelCode
		}
	}
	return 0
}
