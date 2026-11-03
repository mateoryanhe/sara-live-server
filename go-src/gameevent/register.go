package gameevent

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/event"
	"xr-game-server/module/randomnick"
	"xr-game-server/module/promotionattribution"
	userentity "xr-game-server/entity/user"
)

const (
	// RegisterEvent 用户注册成功事件
	RegisterEvent event.Type = "RegisterEvent"
)

// RegisterEventData 用户注册成功事件数据
type RegisterEventData struct {
	UserId                uint64
	RegisteredAt          time.Time
	NickLang              uint8 // 随机昵称语言(0=英文默认)
	Attribution           RegisterAttribution
	PromotionChannelCode  int    // deviceInfo.ch 解析出的渠道编码,0 为无
	ApkPromotionChannel   string // deviceInfo.apkCh,空为无
	PairedPromotionChannelCode int // H5-App 配对成功时的推广URL渠道编码(deviceInfo.ch 或兼容 h5AppCh+apkCh),0 为无
	RegisterAccountChannel     uint   // accounts.channel
	RegisterPromotionPlatform  string // h5/app, 仅按 accounts.channel 落库(不读 Header)
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
	if deviceInfo != nil {
		if c, ok := deviceInfo.PromotionChannelCode(); ok {
			code = c
		}
		if c, ok := deviceInfo.ApkPromotionChannel(); ok {
			apkCh = c
		}
		pairedCode = promotionattribution.PairedPromotionChannelCodeFromDevice(deviceInfo)
	}
	return &RegisterEventData{
		UserId:               userId,
		RegisteredAt:         registeredAt,
		NickLang:             lang,
		Attribution:          registerAttributionFromCtx(ctx, deviceInfo),
		PromotionChannelCode: code,
		ApkPromotionChannel:  apkCh,
		PairedPromotionChannelCode: pairedCode,
	}
}
