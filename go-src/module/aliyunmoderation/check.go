package aliyunmoderation

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/os/gctx"

	"xr-game-server/errercode"
)

// Scene 审核场景，映射阿里云 TextModeration Service
type Scene string

const (
	SceneChat     Scene = "chat"
	SceneNickname Scene = "nickname"
	SceneComment  Scene = "comment"
)

func (s Scene) service(cfg *cfgSnapshot) string {
	switch s {
	case SceneNickname:
		return cfg.NicknameService
	case SceneComment:
		return cfg.CommentService
	default:
		return cfg.ChatService
	}
}

// IsEnabled 是否已开启动态过滤
func IsEnabled() bool {
	return getCfgCache().Enabled
}

// CheckCompliant 检测单段文本；未开启时直接通过
func CheckCompliant(scene Scene, text string) bool {
	return RequireTextCompliant(scene, text) == nil
}

// RequireTextCompliant 校验文本；空串跳过。未开启则跳过；不合规返回 TextSensitiveWord
func RequireTextCompliant(scene Scene, texts ...string) error {
	cfg := getCfgCache()
	if !cfg.Enabled {
		return nil
	}
	if !cfg.ready() {
		return errercode.CreateCode(errercode.AliyunTextModerationCfgInvalid)
	}
	ctx := gctx.New()
	service := scene.service(cfg)
	for _, text := range texts {
		if strings.TrimSpace(text) == "" {
			continue
		}
		ok, err := moderateText(ctx, cfg, service, text)
		if err != nil {
			// 审核服务异常时不放行，避免未审核内容上线
			return errercode.CreateCode(errercode.AliyunTextModerationFailed)
		}
		if !ok {
			return errercode.CreateCode(errercode.TextSensitiveWord)
		}
	}
	return nil
}

// RequireTextCompliantDefault 使用公聊场景检测
func RequireTextCompliantDefault(texts ...string) error {
	return RequireTextCompliant(SceneChat, texts...)
}

// CheckCompliantWithContext 带 context 的检测(供后续扩展)
func CheckCompliantWithContext(ctx context.Context, scene Scene, text string) bool {
	cfg := getCfgCache()
	if !cfg.Enabled || strings.TrimSpace(text) == "" {
		return true
	}
	if !cfg.ready() {
		return false
	}
	ok, err := moderateText(ctx, cfg, scene.service(cfg), text)
	return err == nil && ok
}
