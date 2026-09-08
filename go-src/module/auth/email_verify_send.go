package auth

import (
	"context"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/errercode"
)

var emailCodePattern = regexp.MustCompile(`^\d{6}$`)

// SendEmailVerifyCode 通过 AWS SES 发送验证码邮件.
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

	ec := cfgdao.GetCfEmailCfgCached()
	if !cfgdao.CfEmailEnabled() {
		region, accessKeyId, from := "", "", ""
		enabled := false
		if ec != nil {
			enabled, region, accessKeyId, from = ec.Enabled, ec.Region, ec.AccessKeyId, ec.FromEmail
		}
		xrlog.DetailLog.Warningf(ctx, "email verify send skipped: ses cfg incomplete enabled=%v region=%q accessKeyId=%q from=%q",
			enabled, region, accessKeyId, from)
		return errercode.CreateCode(errercode.SysError)
	}

	subject, body := buildEmailVerifyContent(lang, code)
	client := ses.NewFromConfig(aws.Config{
		Region:      ec.Region,
		Credentials: credentials.NewStaticCredentialsProvider(ec.AccessKeyId, ec.SecretAccessKey, ""),
	})
	_, err := client.SendEmail(ctx, &ses.SendEmailInput{
		Source: aws.String(ec.FromEmail),
		Destination: &types.Destination{
			ToAddresses: []string{toEmail},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
			Body: &types.Body{
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
		},
	})
	if err != nil {
		xrlog.DetailLog.Errorf(ctx, "email verify ses send err to=%s err=%v", toEmail, err)
		return errercode.CreateCode(errercode.SysError)
	}

	markEmailSendSuccess(ctx, toEmail)
	storeEmailVerifyCode(ctx, toEmail, code)
	xrlog.DetailLog.Infof(ctx, "email verify send ok to=%s lang=%s from=%s region=%s", toEmail, normalizeEmailLang(lang), ec.FromEmail, ec.Region)
	return nil
}
