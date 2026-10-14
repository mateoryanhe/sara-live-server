package gameevent

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/event"
	userentity "xr-game-server/entity/user"
	"xr-game-server/module/randomnick"
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
	PromotionChannelCode  int // deviceInfo.ch 解析出的渠道编码,0 为无
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
	if deviceInfo != nil {
		if c, ok := deviceInfo.PromotionChannelCode(); ok {
			code = c
		}
	}
	return &RegisterEventData{
		UserId:               userId,
		RegisteredAt:         registeredAt,
		NickLang:             lang,
		Attribution:          registerAttributionFromCtx(ctx, deviceInfo),
		PromotionChannelCode: code,
	}
}
