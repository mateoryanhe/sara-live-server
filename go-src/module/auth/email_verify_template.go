package auth

import (
	"fmt"
	"strings"
)

// 邮件正文支持的语言(不含中文;与 App 礼物/随机昵称侧一致)
const (
	emailLangEN = "en"
	emailLangES = "es"
	emailLangHI = "hi"
	emailLangPT = "pt"
	emailLangID = "id"
)

type emailVerifyTemplate struct {
	Subject string
	Body    string // 含一个 %s 占位验证码
}

// 写死文案;中文(zh*)及其它未知语言回落英文
var emailVerifyTemplates = map[string]emailVerifyTemplate{
	emailLangEN: {
		Subject: "Your verification code",
		Body: "Your verification code is: %s\n\n" +
			"This code expires in 5 minutes. Do not share it with anyone.\n\n" +
			"If you did not request this code, please ignore this email.",
	},
	emailLangES: {
		Subject: "Tu codigo de verificacion",
		Body: "Tu codigo de verificacion es: %s\n\n" +
			"Este codigo caduca en 5 minutos. No lo compartas con nadie.\n\n" +
			"Si no solicitaste este codigo, ignora este correo.",
	},
	emailLangHI: {
		Subject: "Aapka verification code",
		Body: "Aapka verification code hai: %s\n\n" +
			"Yeh code 5 minute mein expire ho jayega. Kisi ke saath share na karein.\n\n" +
			"Agar aapne yeh code request nahi kiya, is email ko ignore karein.",
	},
	emailLangPT: {
		Subject: "Seu codigo de verificacao",
		Body: "Seu codigo de verificacao e: %s\n\n" +
			"Este codigo expira em 5 minutos. Nao compartilhe com ninguem.\n\n" +
			"Se voce nao solicitou este codigo, ignore este e-mail.",
	},
	emailLangID: {
		Subject: "Kode verifikasi Anda",
		Body: "Kode verifikasi Anda adalah: %s\n\n" +
			"Kode ini kedaluwarsa dalam 5 menit. Jangan bagikan kepada siapa pun.\n\n" +
			"Jika Anda tidak meminta kode ini, abaikan email ini.",
	},
}

func normalizeEmailLang(lang string) string {
	low := strings.ToLower(strings.TrimSpace(lang))
	if i := strings.IndexAny(low, ",;"); i >= 0 {
		low = strings.TrimSpace(low[:i])
	}
	switch {
	case strings.HasPrefix(low, "es"):
		return emailLangES
	case strings.HasPrefix(low, "hi"):
		return emailLangHI
	case strings.HasPrefix(low, "pt"):
		return emailLangPT
	case strings.HasPrefix(low, "id"):
		return emailLangID
	case strings.HasPrefix(low, "en"):
		return emailLangEN
	default:
		// zh-CN / zh-TW 等不支持中文邮件内容,回落英文
		return emailLangEN
	}
}

func buildEmailVerifyContent(lang, code string) (subject, body string) {
	tpl, ok := emailVerifyTemplates[normalizeEmailLang(lang)]
	if !ok {
		tpl = emailVerifyTemplates[emailLangEN]
	}
	return tpl.Subject, fmt.Sprintf(tpl.Body, code)
}
