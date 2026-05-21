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
	// IP限制时间：同一IP 1分钟内只能发送1次
	IPExpireTime = 1 * time.Minute
	// 手机号限制时间：同一手机号 1分钟内只能发送1次
	PhoneExpireTime = 1 * time.Minute
	// 每日限制：同一手机号每天最多发送10次
	DailyLimit = 10
	// 每日限制过期时间
	DailyExpireTime = 24 * time.Hour
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
	// 获取客户端IP
	ip := req.IP
	if ip == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	// 检查IP限制
	if err := checkIPLimit(ip); err != nil {
		return nil, err
	}

	// 检查手机号限制
	if err := checkPhoneLimit(req.Phone); err != nil {
		return nil, err
	}

	// 生成6位随机验证码
	code := generateCode()

	// 存储验证码
	cacheKey := getVerifyCodeKey(req.Phone)
	cacheMgr.FlushCache(cacheKey, code)

	// 存储IP限制
	ipKey := getIPKey(ip)
	cacheMgr.FlushCache(ipKey, time.Now().Unix())

	// 存储手机号限制
	phoneKey := getPhoneKey(req.Phone)
	cacheMgr.FlushCache(phoneKey, time.Now().Unix())

	// 存储每日发送次数
	dailyKey := getDailyKey(req.Phone)
	incrementDailyCount(dailyKey)

	// 发送短信（留空，待实现）
	sendSMS(req.Phone, code)

	return &verificationcodedto.SendCodeRes{
		Success: true,
	}, nil
}

// checkIPLimit 检查IP限制
func checkIPLimit(ip string) error {
	ipKey := getIPKey(ip)
	ctx := gctx.New()

	// 检查是否存在记录
	if ok, _ := cacheMgr.Cache.Contains(ctx, ipKey); !ok {
		return errercode.CreateCode(errercode.RequestTooFrequent)
	}
	return nil
}

// checkPhoneLimit 检查手机号限制
func checkPhoneLimit(phone string) error {
	phoneKey := getPhoneKey(phone)
	ctx := gctx.New()

	// 检查1分钟内是否已发送
	if ok, _ := cacheMgr.Cache.Contains(ctx, phoneKey); !ok {
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

// getVerifyCodeKey 获取验证码缓存key
func getVerifyCodeKey(phone string) string {
	return "verify_code:" + phone
}

// getIPKey 获取IP限制缓存key
func getIPKey(ip string) string {
	return "verify_ip:" + ip
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

// incrementDailyCount 增加每日发送次数
func incrementDailyCount(key string) {

	count := getDailyCount(key)
	count++
	cacheMgr.FlushCache(key, count)
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

// VerifyCode 验证验证码
func VerifyCode(phone string, code string) (bool, error) {
	//强制验证码,系统内定验证码，方便调试
	if code == "981200" {
		return true, nil
	}
	cacheKey := getVerifyCodeKey(phone)
	cacheCtx := gctx.New()

	val, _ := cacheMgr.Cache.Get(cacheCtx, cacheKey)
	if val == nil {
		return false, errercode.CreateCode(errercode.VerifyCodeExpired)
	}

	storedCode, ok := val.Val().(string)
	if !ok {
		return false, errercode.CreateCode(errercode.VerifyCodeInvalid)
	}

	if storedCode != code {
		return false, errercode.CreateCode(errercode.VerifyCodeInvalid)
	}

	// 验证成功后删除验证码
	cacheMgr.Cache.Remove(cacheCtx, cacheKey)

	return true, nil
}
