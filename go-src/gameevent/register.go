package gameevent

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/event"
	userentity "xr-game-server/entity/user"
	"xr-game-server/module/promotionattribution"
	"xr-game-server/module/randomnick"
)

const (
	// RegisterEvent 用户注册成功事件
	RegisterEvent event.Type = "RegisterEvent"
)

// RegisterEventData 用户注册成功事件数据
type RegisterEventData struct {
	UserId                     uint64
	RegisteredAt               time.Time
	NickLang                   uint8 // 随机昵称语言(0=英文默认)
	Attribution                RegisterAttribution
	PromotionChannelCode       int    // deviceInfo.ch 解析出的渠道编码,0 为无
	ApkPromotionChannel        string // deviceInfo.apkCh,空为无
	PairedPromotionChannelCode int    // H5-App 配对成功时的推广URL渠道编码(deviceInfo.ch 或兼容 h5AppCh+apkCh),0 为无
	RegisterAccountChannel     uint   // accounts.channel
	RegisterPromotionPlatform  string // h5/app, 优先 X-Platform(web→h5,android/ios→app),否则 accounts.channel
	// RegisterUrlAttr* 来自 deviceInfo.registerAttr; HasRegisterUrlAttr=false 表示客户端未上报,不上库/不上报 accept
	HasRegisterUrlAttr bool
	RegisterUrlAttrType     string
	RegisterUrlAttrCh       string
	RegisterUrlAttrTjUserId string
	DeviceModel             string
	OsVersion               string
}

func NewRegisterEventData(userId uint64, registeredAt time.Time) *RegisterEventData {
	return &RegisterEventData{
		UserId:       userId,
		RegisteredAt: registeredAt,
		NickLang:     randomnick.DefaultLang,
	}
}

func NewRegisterEventDataFromCtx(ctx context.Context, userId uint64, registeredAt time.Time) *RegisterEventData {
	return NewRegisterEventDataFromCtxWithDevice(ctx, userId, registeredAt, nil)
}

func NewRegisterEventDataFromCtxWithDevice(ctx context.Context, userId uint64, registeredAt time.Time, deviceInfo *userentity.DeviceInfo) *RegisterEventData {
	lang := randomnick.DefaultLang
	if r := g.RequestFromCtx(ctx); r != nil {
		lang = randomnick.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	}
	code := 0
	apkCh := ""
	pairedCode := 0
	deviceModel := ""
	osVersion := ""
	attrType, attrCh, attrTj := "", "", ""
	hasAttr := false
	if deviceInfo != nil {
		if c, ok := deviceInfo.PromotionChannelCode(); ok {
			code = c
		}
		if c, ok := deviceInfo.ApkPromotionChannel(); ok {
			apkCh = c
		}
		pairedCode = promotionattribution.PairedPromotionChannelCodeFromDevice(deviceInfo)
		deviceModel = strings.TrimSpace(deviceInfo.DeviceModel)
		osVersion = strings.TrimSpace(deviceInfo.OsVersion)
		if t, c, tj, ok := deviceInfo.RegisterAttr.Normalize(); ok {
			hasAttr = true
			attrType, attrCh, attrTj = t, c, tj
		}
	}
	return &RegisterEventData{
		UserId:                     userId,
		RegisteredAt:               registeredAt,
		NickLang:                   lang,
		Attribution:                registerAttributionFromCtx(ctx, deviceInfo),
		PromotionChannelCode:       code,
		ApkPromotionChannel:        apkCh,
		PairedPromotionChannelCode: pairedCode,
		HasRegisterUrlAttr:         hasAttr,
		RegisterUrlAttrType:        attrType,
		RegisterUrlAttrCh:          attrCh,
		RegisterUrlAttrTjUserId:    attrTj,
		DeviceModel:                deviceModel,
		OsVersion:                  osVersion,
	}
}
