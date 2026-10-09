package auth

import (
	"context"
	"regexp"
	"strings"

	"xr-game-server/core/xrlog"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/errercode"
)

var emailCodePattern = regexp.MustCompile(`^\d{6}$`)

// SendEmailVerifyCode 通过 Amazon SES SMTP 发送验证码邮件.
// 发信成功后写入验证码缓存(5分钟,见 email_verify_code.go);生成由 App 接口 SendEmailCode 完成.
// 限流(进程内 gcache):同邮箱 1 分钟冷却 + 每日最多 10 次(本地 0 点过期).
// code 必须为 6 位数字; lang 支持 en/es/hi/pt/id,中文及其它回落 en.
// 参数取自 CMS 配置表 cf_email_cfgs(内存缓存);发件身份/域须已在 SES 验证.
func SendEmailVerifyCode(ctx context.Context, toEmail, code, lang string) error {
	toEmail = normalizeEmailKey(toEmail)
	code = strings.TrimSpace(code)
	if toEmail == "" || !emailCodePattern.MatchString(code) {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if err := checkEmailSendLimit(ctx, toEmail); err != nil {
		return err
	}
	if err := sendEmailVerifyCodeOnce(ctx, toEmail, code, lang); err != nil {
		return err
	}
	markEmailSendSuccess(ctx, toEmail)
	storeEmailVerifyCode(ctx, toEmail, code)
	return nil
}

// SendEmailVerifyCodeForTest CMS 测试发信:生成 6 位码并 SMTP 发送,不占 App 限流、不写入登录验证码缓存.
func SendEmailVerifyCodeForTest(ctx context.Context, toEmail, lang string) error {
	toEmail = normalizeEmailKey(toEmail)
	if toEmail == "" {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	code, err := generateEmailVerifyCode()
	if err != nil {
		return err
	}
	return sendEmailVerifyCodeOnce(ctx, toEmail, code, lang)
}

func sendEmailVerifyCodeOnce(ctx context.Context, toEmail, code, lang string) error {
	ec := cfgdao.GetCfEmailCfgCached()
	if !cfgdao.CfEmailEnabled() {
		host, username, from := "", "", ""
		enabled := false
		port := 0
		if ec != nil {
			enabled = ec.Enabled
			host = ec.SmtpHost
			port = ec.SmtpPort
			username = ec.SmtpUsername
			from = ec.FromEmail
		}
		xrlog.DetailLog.Warningf(ctx, "email verify send skipped: smtp cfg incomplete enabled=%v host=%q port=%d username=%q from=%q",
			enabled, host, port, username, from)
		return errercode.CreateCode(errercode.SysError)
	}

	subject, body := buildEmailVerifyContent(lang, code)
	if err := sendEmailViaSMTP(ctx, ec, toEmail, subject, body); err != nil {
		xrlog.DetailLog.Errorf(ctx, "email verify smtp send err to=%s err=%v", toEmail, err)
		return errercode.CreateCode(errercode.SysError)
	}
	host, port, _, _, from, _ := cfEmailSMTPSettings(ec)
	xrlog.DetailLog.Infof(ctx, "email verify send ok to=%s lang=%s from=%s smtp=%s:%d", toEmail, normalizeEmailLang(lang), from, host, port)
	return nil
}
