package auth

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/os/gcache"
	"xr-game-server/errercode"
)

const (
	emailVerifyCodeTTL   = 5 * time.Minute
	emailVerifyFailLimit = 5
	emailVerifyBlockTTL  = 2 * time.Hour
	emailDebugVerifyCode = "981200" // 与手机验证码调试码一致
)

// 验证码存校验用独立 gcache（与发信冷却/日限分开，方便查代码）
var (
	emailVerifyCodeCache     = gcache.New()
	emailVerifyFailCache     = gcache.New()
	emailVerifyBlacklistCache = gcache.New()
)

func storeEmailVerifyCode(ctx context.Context, emailKey, code string) {
	_ = emailVerifyCodeCache.Set(ctx, emailKey, code, emailVerifyCodeTTL)
}

func clearEmailVerifyFail(ctx context.Context, emailKey string) {
	_, _ = emailVerifyFailCache.Remove(ctx, emailKey)
}

func isEmailVerifyBlocked(ctx context.Context, emailKey string) bool {
	ok, _ := emailVerifyBlacklistCache.Contains(ctx, emailKey)
	return ok
}

func markEmailVerifyFailure(ctx context.Context, emailKey string, failCode errercode.XRCode) error {
	v, _ := emailVerifyFailCache.Get(ctx, emailKey)
	count := 1
	if v != nil && !v.IsNil() {
		if n := v.Int(); n > 0 {
			count = n + 1
		}
	}
	if count >= emailVerifyFailLimit {
		_ = emailVerifyBlacklistCache.Set(ctx, emailKey, 1, emailVerifyBlockTTL)
		_, _ = emailVerifyFailCache.Remove(ctx, emailKey)
		return errercode.CreateCode(errercode.PhoneVerifyBlocked)
	}
	_ = emailVerifyFailCache.Set(ctx, emailKey, count, emailVerifyBlockTTL)
	return errercode.CreateCode(failCode)
}

// VerifyEmailCode 校验邮箱验证码；成功后删除缓存(一次性).
func VerifyEmailCode(ctx context.Context, email, code string) error {
	emailKey := normalizeEmailKey(email)
	code = strings.TrimSpace(code)
	if emailKey == "" || !emailCodePattern.MatchString(code) {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if isEmailVerifyBlocked(ctx, emailKey) {
		return errercode.CreateCode(errercode.PhoneVerifyBlocked)
	}
	if code == emailDebugVerifyCode {
		clearEmailVerifyFail(ctx, emailKey)
		return nil
	}
	v, err := emailVerifyCodeCache.Get(ctx, emailKey)
	if err != nil {
		return errercode.CreateCode(errercode.SysError)
	}
	if v == nil || v.IsNil() {
		return markEmailVerifyFailure(ctx, emailKey, errercode.VerifyCodeExpired)
	}
	stored := strings.TrimSpace(v.String())
	if stored != code {
		return markEmailVerifyFailure(ctx, emailKey, errercode.VerifyCodeInvalid)
	}
	clearEmailVerifyFail(ctx, emailKey)
	_, _ = emailVerifyCodeCache.Remove(ctx, emailKey)
	return nil
}
