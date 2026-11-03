package promotionattribution

import (
	"regexp"
	"strings"

	userentity "xr-game-server/entity/user"
	"xr-game-server/dao/apkpromotionchanneldao"
	"xr-game-server/dao/promotionchannelurldao"
)

var h5AppChannelCodePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)

// NormalizeH5AppChannelCode 与 APK 渠道码相同字符集.
func NormalizeH5AppChannelCode(raw string) (string, bool) {
	code := strings.TrimSpace(raw)
	if code == "" || !h5AppChannelCodePattern.MatchString(code) {
		return "", false
	}
	return code, true
}

// ResolveH5AppPairByPromotionChannelCode App/H5 仅上报 deviceInfo.ch(推广 URL 整数编码) 时,
// 按 CMS 中 URL.h5_app_channel_code 与 APK 行配对, 返回 APK channel_code.
func ResolveH5AppPairByPromotionChannelCode(promoCode int) (apkChannelCode string, ok bool) {
	if promoCode <= 0 {
		return "", false
	}
	urlRow := promotionchannelurldao.GetByChannelCode(promoCode)
	if urlRow == nil || urlRow.ChannelCode <= 0 {
		return "", false
	}
	h5App := strings.TrimSpace(urlRow.H5AppChannelCode)
	if h5App == "" {
		return "", false
	}
	apkRow := apkpromotionchanneldao.GetFirstByH5AppChannelCode(h5App)
	if apkRow == nil {
		return "", false
	}
	if strings.TrimSpace(apkRow.H5AppChannelCode) != h5App {
		return "", false
	}
	code := strings.TrimSpace(apkRow.ChannelCode)
	if code == "" {
		return "", false
	}
	return code, true
}

// NormalizeApkPromotionChannelCode deviceInfo.apkCh 合法格式(与 CMS 是否预建无关).
func NormalizeApkPromotionChannelCode(raw string) (string, bool) {
	return NormalizeH5AppChannelCode(raw)
}

// ResolvePromotionURLChannelCodeByApkCh App 上报 apkCh 时, 与推广 URL 列表 h5_app_channel_code 匹配即得 URL 整数渠道.
func ResolvePromotionURLChannelCodeByApkCh(apkCh string) (promotionChannelCode int, ok bool) {
	apkCh, ok = NormalizeApkPromotionChannelCode(apkCh)
	if !ok {
		return 0, false
	}
	urlRow := promotionchannelurldao.GetByH5AppChannelCode(apkCh)
	if urlRow == nil || urlRow.ChannelCode <= 0 {
		return 0, false
	}
	if strings.TrimSpace(urlRow.H5AppChannelCode) != apkCh {
		return 0, false
	}
	return urlRow.ChannelCode, true
}

// ResolvePromotionChannelCodeByApkChannelCode App 注册 deviceInfo.apkCh → 推广 URL 整数 channel_code(先 URL 匹配, 再 APK CMS 配对).
func ResolvePromotionChannelCodeByApkChannelCode(apkCh string) (promotionChannelCode int, ok bool) {
	if c, ok := ResolvePromotionURLChannelCodeByApkCh(apkCh); ok {
		return c, true
	}
	apkCh = strings.TrimSpace(apkCh)
	if apkCh == "" {
		return 0, false
	}
	apkRow := apkpromotionchanneldao.GetByChannelCode(apkCh)
	if apkRow == nil {
		return 0, false
	}
	h5App := strings.TrimSpace(apkRow.H5AppChannelCode)
	if h5App == "" {
		return 0, false
	}
	return ValidatePair(h5App, apkCh)
}

// PairedPromotionChannelCodeFromDevice H5-App 推广日统计用配对渠道码.
// H5: deviceInfo.ch; App: deviceInfo.apkCh(经 CMS 反查 URL 渠道); 兼容 h5AppCh+apkCh.
func PairedPromotionChannelCodeFromDevice(info *userentity.DeviceInfo) int {
	if info == nil {
		return 0
	}
	if promo, ok := info.PromotionChannelCode(); ok {
		if _, paired := ResolveH5AppPairByPromotionChannelCode(promo); paired {
			return promo
		}
	}
	if apk, ok := info.ApkPromotionChannel(); ok {
		if c, ok := ResolvePromotionChannelCodeByApkChannelCode(apk); ok {
			return c
		}
	}
	if h5, ok1 := info.H5AppChannel(); ok1 {
		if apk, ok2 := info.ApkPromotionChannel(); ok2 {
			if c, ok := ValidatePair(h5, apk); ok {
				return c
			}
		}
	}
	return 0
}

// ValidatePair App 显式上报 h5AppCh+apkCh 时须与 CMS 配置一致(兼容旧客户端).
// 返回推广 URL 列表中的整数 channelCode(用于日统计维度).
func ValidatePair(h5AppCh, apkCh string) (promotionChannelCode int, ok bool) {
	h5AppCh = strings.TrimSpace(h5AppCh)
	apkCh = strings.TrimSpace(apkCh)
	if h5AppCh == "" || apkCh == "" {
		return 0, false
	}
	urlRow := promotionchannelurldao.GetByH5AppChannelCode(h5AppCh)
	if urlRow == nil || urlRow.ChannelCode <= 0 {
		return 0, false
	}
	if strings.TrimSpace(urlRow.H5AppChannelCode) != h5AppCh {
		return 0, false
	}
	apkRow := apkpromotionchanneldao.GetByChannelCode(apkCh)
	if apkRow == nil {
		return 0, false
	}
	if strings.TrimSpace(apkRow.H5AppChannelCode) != h5AppCh {
		return 0, false
	}
	return urlRow.ChannelCode, true
}

// UserHasValidPair 已绑定用户 H5-App 日统计(下单/充值): apkCh 与 URL 渠道码一致即可, 不强制 APK CMS 行.
func UserHasValidPair(promotionChannelCode int, apkCh string) bool {
	if promotionChannelCode <= 0 {
		return false
	}
	apkCh = strings.TrimSpace(apkCh)
	if apkCh == "" {
		var ok bool
		apkCh, ok = ResolveH5AppPairByPromotionChannelCode(promotionChannelCode)
		if !ok {
			return false
		}
	}
	if c, ok := ResolvePromotionURLChannelCodeByApkCh(apkCh); ok && c == promotionChannelCode {
		return true
	}
	urlRow := promotionchannelurldao.GetByChannelCode(promotionChannelCode)
	if urlRow == nil {
		return false
	}
	h5App := strings.TrimSpace(urlRow.H5AppChannelCode)
	if h5App == "" {
		return false
	}
	_, ok := ValidatePair(h5App, apkCh)
	return ok
}
