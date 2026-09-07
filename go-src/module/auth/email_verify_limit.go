package auth

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/os/gcache"
	"xr-game-server/errercode"
)

const (
	// 对齐手机验证码：同一邮箱 1 分钟内只能发 1 次
	emailSendCooldownTTL = 1 * time.Minute
	emailSendDailyLimit  = 10
)

// 两个独立 gcache，方便查代码：发信冷却 / 每日次数（到本地 0 点过期，不入库）
var (
	emailSendCooldownCache = gcache.New()
	emailSendDailyCache    = gcache.New()
)

func normalizeEmailKey(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func ttlUntilNextLocalMidnight() time.Duration {
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	d := next.Sub(now)
	if d < time.Second {
		return time.Second
	}
	return d
}

// checkEmailSendLimit 发信前检查冷却与每日上限（不含验证码存储）
func checkEmailSendLimit(ctx context.Context, emailKey string) error {
	ok, err := emailSendCooldownCache.Contains(ctx, emailKey)
	if err != nil {
		return errercode.CreateCode(errercode.SysError)
	}
	if ok {
		return errercode.CreateCode(errercode.RequestTooFrequent)
	}
	count, err := getEmailSendDailyCount(ctx, emailKey)
	if err != nil {
		return errercode.CreateCode(errercode.SysError)
	}
	if count >= emailSendDailyLimit {
		return errercode.CreateCode(errercode.DailyLimitExceeded)
	}
	return nil
}

func getEmailSendDailyCount(ctx context.Context, emailKey string) (int, error) {
	v, err := emailSendDailyCache.Get(ctx, emailKey)
	if err != nil {
		return 0, err
	}
	if v == nil || v.IsNil() {
		return 0, nil
	}
	return v.Int(), nil
}

// markEmailSendSuccess 仅在 CF 发信成功后占用冷却并累加当日次数
func markEmailSendSuccess(ctx context.Context, emailKey string) {
	_ = emailSendCooldownCache.Set(ctx, emailKey, 1, emailSendCooldownTTL)
	count, _ := getEmailSendDailyCount(ctx, emailKey)
	_ = emailSendDailyCache.Set(ctx, emailKey, count+1, ttlUntilNextLocalMidnight())
}
