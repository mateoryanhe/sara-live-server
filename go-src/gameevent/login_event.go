package gameevent

import (
	"context"

	"xr-game-server/constants/lang"
	"xr-game-server/core/event"
	"xr-game-server/core/httpserver"
)

// LoginEvent App 用户登录成功(含 getUserInfo 拉取本人信息).
const LoginEvent event.Type = "LoginEvent"

// LoginEventData 登录事件载荷.
type LoginEventData struct {
	UserId         uint64
	ClientPlatform string // 请求头 X-Platform 规范化: android / ios / web
	ClientLang     lang.Lang // X-Language / Accept-Language
}

func NewLoginEventDataFromCtx(ctx context.Context, userId uint64) *LoginEventData {
	return &LoginEventData{
		UserId:         userId,
		ClientPlatform: StatClientPlatformFromCtx(ctx),
		ClientLang:     httpserver.GetLangFromContext(ctx),
	}
}

// PublishLoginEvent 发布登录事件(ClientPlatform 在发布时从 ctx 读取).
func PublishLoginEvent(ctx context.Context, userId uint64) {
	if userId == 0 {
		return
	}
	event.Pub(LoginEvent, NewLoginEventDataFromCtx(ctx, userId))
}

// ParseLoginEventUserId 兼容历史 uint64 载荷与新 LoginEventData.
func ParseLoginEventUserId(data any) (uint64, bool) {
	switch v := data.(type) {
	case uint64:
		return v, v > 0
	case *LoginEventData:
		if v == nil || v.UserId == 0 {
			return 0, false
		}
		return v.UserId, true
	default:
		return 0, false
	}
}