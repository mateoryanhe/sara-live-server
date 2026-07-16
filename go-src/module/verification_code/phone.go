package verification_code

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gctx"
	"math/rand"
	"time"
	"xr-game-server/core/cache"
	"xr-game-server/dto/verificationcodedto"
	"xr-game-server/errercode"
)

const (
	// 验证码有效期
	CodeExpireTime = 5 * time.Minute
	// 手机号限制时间：同一手机号 1分钟内只能发送1次
	PhoneExpireTime = 1 * time.Minute
	// 每日限制：同一手机号每天最多发送10次
	DailyLimit = 10
	// 每日限制过期时间
	DailyExpireTime = 24 * time.Hour
	// 验证码连续失败上限
	VerifyFailLimit = 5
	// 手机号校验失败拉黑时间
	PhoneBlacklistTime = 2 * time.Hour
)

var (
	cacheMgr *cache.CacheMgr
)

// Init 初始化验证码模块
func Init() {
	cacheMgr = cache.NewCacheMgr()
}

// SendCode 发送验证码
func SendCode(ctx context.Context, req *verificationcodedto.SendCodeReq) (*verificationcodedto.SendCodeRes, error) {
	_ = ctx

	// 检查手机号限制
	if err := checkPhoneLimit(req.Phone); err != nil {
		return nil, err
	}

	// 生成6位随机验证码
	code := generateCode()

	// 存储验证码
	cacheKey := getVerifyCodeKey(req.Phone)
	setCacheWithTTL(cacheKey, code, CodeExpireTime)

	// 存储手机号限制
	phoneKey := getPhoneKey(req.Phone)
	setCacheWithTTL(phoneKey, time.Now().Unix(), PhoneExpireTime)

	// 存储每日发送次数
	dailyKey := getDailyKey(req.Phone)
	incrementDailyCount(dailyKey)

	// 发送短信（留空，待实现）
	sendSMS(req.Phone, code)

	return &verificationcodedto.SendCodeRes{
		Success: true,
	}, nil
}

// checkPhoneLimit 检查手机号限制
func checkPhoneLimit(phone string) error {
	phoneKey := getPhoneKey(phone)
	ctx := gctx.New()

	// 1分钟内已发送则拒绝
	if ok, _ := cacheMgr.Cache.Contains(ctx, phoneKey); ok {
		return errercode.CreateCode(errercode.RequestTooFrequent)
	}

	// 检查每日限制
	dailyKey := getDailyKey(phone)
	count := getDailyCount(dailyKey)
	if count >= DailyLimit {
		return errercode.CreateCode(errercode.DailyLimitExceeded)
	}

	return nil
}

// generateCode 生成6位随机验证码
func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000
	return fmt.Sprintf("%06d", code)
}

// sendSMS 发送短信（留空，待实现）
func sendSMS(phone, code string) {
	// TODO: 实现短信发送逻辑
	// 目前短信运营商还未确定，留空
}

func setCacheWithTTL(key interface{}, data any, ttl time.Duration) {
	_ = cacheMgr.Cache.Set(gctx.New(), key, data, ttl)
}

// getVerifyCodeKey 获取验证码缓存key
func getVerifyCodeKey(phone string) string {
	return "verify_code:" + phone
}

// getPhoneKey 获取手机号限制缓存key
func getPhoneKey(phone string) string {
	return "verify_phone:" + phone
}

// getDailyKey 获取每日限制缓存key
func getDailyKey(phone string) string {
	date := time.Now().Format("2006-01-02")
	return "verify_daily:" + phone + ":" + date
}

func getVerifyFailCountKey(phone string) string {
	return "verify_fail_count:" + phone
}

func getPhoneBlacklistKey(phone string) string {
	return "verify_blacklist:" + phone
}

// incrementDailyCount 增加每日发送次数
func incrementDailyCount(key string) {
	count := getDailyCount(key)
	count++
	setCacheWithTTL(key, count, DailyExpireTime)
}

// getDailyCount 获取每日发送次数
func getDailyCount(key string) int {
	ctx := gctx.New()
	val, _ := cacheMgr.Cache.Get(ctx, key)
	if val == nil {
		return 0
	}
	count, ok := val.Val().(int)
	if !ok {
		return 0
	}
	return count
}

func isPhoneBlacklisted(phone string) bool {
	ctx := gctx.New()
	ok, _ := cacheMgr.Cache.Contains(ctx, getPhoneBlacklistKey(phone))
	return ok
}

func clearVerifyFailCount(phone string) {
	_, _ = cacheMgr.Cache.Remove(gctx.New(), getVerifyFailCountKey(phone))
}

func onVerifyFailure(phone string, failCode errercode.XRCode) (bool, error) {
	if blockedErr := markVerifyFailure(phone); blockedErr != nil {
		return false, blockedErr
	}
	return false, errercode.CreateCode(failCode)
}

func markVerifyFailure(phone string) error {
	ctx := gctx.New()
	key := getVerifyFailCountKey(phone)
	count := 1

	val, _ := cacheMgr.Cache.Get(ctx, key)
	if val != nil {
		if current, ok := val.Val().(int); ok && current > 0 {
			count = current + 1
		}
	}

	if count >= VerifyFailLimit {
		setCacheWithTTL(getPhoneBlacklistKey(phone), time.Now().Unix(), PhoneBlacklistTime)
		_, _ = cacheMgr.Cache.Remove(ctx, key)
		return errercode.CreateCode(errercode.PhoneVerifyBlocked)
	}

	setCacheWithTTL(key, count, PhoneBlacklistTime)
	return nil
}

// VerifyCode 验证验证码
func VerifyCode(phone string, code string) (bool, error) {
	if isPhoneBlacklisted(phone) {
		return false, errercode.CreateCode(errercode.PhoneVerifyBlocked)
	}

	//强制验证码,系统内定验证码，方便调试
	if code == "981200" {
		clearVerifyFailCount(phone)
		return true, nil
	}

	cacheKey := getVerifyCodeKey(phone)
	cacheCtx := gctx.New()

	val, _ := cacheMgr.Cache.Get(cacheCtx, cacheKey)
	if val == nil {
		return onVerifyFailure(phone, errercode.VerifyCodeExpired)
	}

	storedCode, ok := val.Val().(string)
	if !ok {
		return onVerifyFailure(phone, errercode.VerifyCodeInvalid)
	}

	if storedCode != code {
		return onVerifyFailure(phone, errercode.VerifyCodeInvalid)
	}

	clearVerifyFailCount(phone)
	_, _ = cacheMgr.Cache.Remove(cacheCtx, cacheKey)

	return true, nil
}
