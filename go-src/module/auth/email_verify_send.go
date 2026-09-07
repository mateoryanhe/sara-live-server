package auth

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/errercode"
)

var emailCodePattern = regexp.MustCompile(`^\d{6}$`)

const cfEmailSendURLFmt = "https://api.cloudflare.com/client/v4/accounts/%s/email/sending/send"

// SendEmailVerifyCode 通过 Cloudflare Email Sending REST API 发送验证码邮件.
// 发信成功后写入验证码缓存(5分钟,见 email_verify_code.go);生成由 App 接口 SendEmailCode 完成.
// 限流(进程内 gcache):同邮箱 1 分钟冷却 + 每日最多 10 次(本地 0 点过期).
// code 必须为 6 位数字; lang 支持 en/es/hi/pt/id,中文及其它回落 en.
// 参数取自 CMS 配置表 cf_email_cfgs(内存缓存);发件域需已在 CF Email Service Onboard.
func SendEmailVerifyCode(ctx context.Context, toEmail, code, lang string) error {
	toEmail = normalizeEmailKey(toEmail)
	code = strings.TrimSpace(code)
	if toEmail == "" || !emailCodePattern.MatchString(code) {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if err := checkEmailSendLimit(ctx, toEmail); err != nil {
		return err
	}

	ec := cfgdao.GetCfEmailCfgCached()
	if !cfgdao.CfEmailEnabled() {
		from, accountId := "", ""
		enabled := false
		if ec != nil {
			enabled, from, accountId = ec.Enabled, ec.FromEmail, ec.AccountId
		}
		xrlog.DetailLog.Warningf(ctx, "email verify send skipped: cf email cfg incomplete enabled=%v from=%q accountId=%q",
			enabled, from, accountId)
		return errercode.CreateCode(errercode.SysError)
	}

	subject, body := buildEmailVerifyContent(lang, code)
	payload := g.Map{
		"to":      toEmail,
		"from":    ec.FromEmail,
		"subject": subject,
		"text":    body,
	}
	url := fmt.Sprintf(cfEmailSendURLFmt, ec.AccountId)
	client := g.Client().Clone()
	client.SetHeader("Authorization", "Bearer "+ec.ApiToken)
	client.SetHeader("Content-Type", "application/json")
	resp, err := client.ContentJson().Post(ctx, url, payload)
	if err != nil {
		xrlog.DetailLog.Errorf(ctx, "email verify cf send http err to=%s err=%v", toEmail, err)
		return errercode.CreateCode(errercode.SysError)
	}
	defer resp.Close()
	raw := resp.ReadAllString()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		xrlog.DetailLog.Errorf(ctx, "email verify cf send http status=%d to=%s body=%s", resp.StatusCode, toEmail, truncateEmailLog(raw, 500))
		return errercode.CreateCode(errercode.SysError)
	}
	j, err := gjson.DecodeToJson(raw)
	if err != nil {
		xrlog.DetailLog.Errorf(ctx, "email verify cf send parse err to=%s body=%s err=%v", toEmail, truncateEmailLog(raw, 500), err)
		return errercode.CreateCode(errercode.SysError)
	}
	if !j.Get("success").Bool() {
		xrlog.DetailLog.Errorf(ctx, "email verify cf send failed to=%s body=%s", toEmail, truncateEmailLog(raw, 500))
		return errercode.CreateCode(errercode.SysError)
	}
	markEmailSendSuccess(ctx, toEmail)
	storeEmailVerifyCode(ctx, toEmail, code)
	xrlog.DetailLog.Infof(ctx, "email verify send ok to=%s lang=%s from=%s", toEmail, normalizeEmailLang(lang), ec.FromEmail)
	return nil
}

func truncateEmailLog(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
