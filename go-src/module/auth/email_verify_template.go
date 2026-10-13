package auth

import (
	"fmt"
	"html"
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
	Subject      string
	Title        string
	Intro        string
	ExpiryNote   string
	SecurityNote string
	IgnoreNote   string
	Footer       string
}

// 写死文案;中文(zh*)及其它未知语言回落英文
var emailVerifyTemplates = map[string]emailVerifyTemplate{
	emailLangEN: {
		Subject:      "Your verification code",
		Title:        "Verify your email",
		Intro:        "Enter this verification code in the app to continue. The code is valid for a single sign-in attempt.",
		ExpiryNote:   "This code expires in 5 minutes.",
		SecurityNote: "Never share this code with anyone. Our team will never ask for it.",
		IgnoreNote:   "If you did not request this code, you can safely ignore this email.",
		Footer:       "This is an automated message. Please do not reply.",
	},
	emailLangES: {
		Subject:      "Tu codigo de verificacion",
		Title:        "Verifica tu correo",
		Intro:        "Introduce este codigo en la aplicacion para continuar. El codigo es valido para un solo intento.",
		ExpiryNote:   "Este codigo caduca en 5 minutos.",
		SecurityNote: "No compartas este codigo con nadie. Nuestro equipo nunca te lo pedira.",
		IgnoreNote:   "Si no solicitaste este codigo, puedes ignorar este correo.",
		Footer:       "Este es un mensaje automatico. No respondas a este correo.",
	},
	emailLangHI: {
		Subject:      "Aapka verification code",
		Title:        "Apna email verify karein",
		Intro:        "Aage badhne ke liye app mein yeh verification code darj karein. Yeh code ek baar ke login ke liye hai.",
		ExpiryNote:   "Yeh code 5 minute mein expire ho jayega.",
		SecurityNote: "Is code ko kisi ke saath share na karein. Hamari team kabhi yeh nahi maangegi.",
		IgnoreNote:   "Agar aapne yeh code request nahi kiya, is email ko ignore kar sakte hain.",
		Footer:       "Yeh ek automatic message hai. Is par reply na karein.",
	},
	emailLangPT: {
		Subject:      "Seu codigo de verificacao",
		Title:        "Verifique seu e-mail",
		Intro:        "Digite este codigo no aplicativo para continuar. O codigo vale para uma unica tentativa.",
		ExpiryNote:   "Este codigo expira em 5 minutos.",
		SecurityNote: "Nunca compartilhe este codigo. Nossa equipe nunca pedira esse codigo.",
		IgnoreNote:   "Se voce nao solicitou este codigo, ignore este e-mail.",
		Footer:       "Esta e uma mensagem automatica. Nao responda.",
	},
	emailLangID: {
		Subject:      "Kode verifikasi Anda",
		Title:        "Verifikasi email Anda",
		Intro:        "Masukkan kode verifikasi ini di aplikasi untuk melanjutkan. Kode hanya berlaku untuk satu kali percobaan.",
		ExpiryNote:   "Kode ini kedaluwarsa dalam 5 menit.",
		SecurityNote: "Jangan bagikan kode ini kepada siapa pun. Tim kami tidak akan pernah memintanya.",
		IgnoreNote:   "Jika Anda tidak meminta kode ini, abaikan email ini.",
		Footer:       "Pesan otomatis. Mohon jangan membalas email ini.",
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

func emailVerifyTemplateForLang(lang string) emailVerifyTemplate {
	tpl, ok := emailVerifyTemplates[normalizeEmailLang(lang)]
	if !ok {
		return emailVerifyTemplates[emailLangEN]
	}
	return tpl
}

func buildEmailVerifyContent(lang, code string) (subject, plainBody, htmlBody string) {
	tpl := emailVerifyTemplateForLang(lang)
	subject = tpl.Subject
	plainBody = buildEmailVerifyPlainBody(tpl, code)
	htmlBody = buildEmailVerifyHTMLBody(tpl, code)
	return subject, plainBody, htmlBody
}

func buildEmailVerifyPlainBody(tpl emailVerifyTemplate, code string) string {
	var b strings.Builder
	b.WriteString(tpl.Title)
	b.WriteString("\n\n")
	b.WriteString(tpl.Intro)
	b.WriteString("\n\n")
	b.WriteString(code)
	b.WriteString("\n\n")
	b.WriteString(tpl.ExpiryNote)
	b.WriteString("\n")
	b.WriteString(tpl.SecurityNote)
	b.WriteString("\n\n")
	b.WriteString(tpl.IgnoreNote)
	b.WriteString("\n\n")
	b.WriteString("—\n")
	b.WriteString(tpl.Footer)
	return b.String()
}

func buildEmailVerifyHTMLBody(tpl emailVerifyTemplate, code string) string {
	esc := func(s string) string { return html.EscapeString(s) }
	codeEsc := esc(code)
	// 内联样式,兼容常见邮件客户端;配色参考 GitHub 通知邮件(无图片)
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"></head>
<body style="margin:0;padding:0;background-color:#f6f8fa;">
<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="background-color:#f6f8fa;">
<tr><td align="center" style="padding:32px 16px;">
<table role="presentation" width="100%%" cellspacing="0" cellpadding="0" style="max-width:480px;background-color:#ffffff;border:1px solid #d0d7de;border-radius:8px;overflow:hidden;">
<tr><td style="padding:24px 28px 20px;border-bottom:1px solid #d0d7de;background:linear-gradient(180deg,#ffffff 0%%,#f6f8fa 100%%);">
<p style="margin:0;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Helvetica,Arial,sans-serif;font-size:20px;font-weight:600;color:#24292f;line-height:1.3;">%s</p>
</td></tr>
<tr><td style="padding:28px 28px 8px;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Helvetica,Arial,sans-serif;">
<p style="margin:0 0 20px;font-size:14px;line-height:1.6;color:#57606a;">%s</p>
<table role="presentation" width="100%%" cellspacing="0" cellpadding="0">
<tr><td align="center" style="padding:8px 0 24px;">
<div style="display:inline-block;padding:18px 28px;background-color:#ddf4ff;border:1px solid #54aeff;border-radius:8px;">
<span style="font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;font-size:32px;font-weight:700;letter-spacing:0.35em;color:#0969da;">%s</span>
</div>
</td></tr>
</table>
<p style="margin:0 0 8px;font-size:13px;line-height:1.5;color:#24292f;"><strong style="color:#1a7f37;">&#9679;</strong> <span style="color:#57606a;">%s</span></p>
<p style="margin:0 0 20px;font-size:13px;line-height:1.5;color:#57606a;">%s</p>
<p style="margin:0;padding:12px 14px;font-size:13px;line-height:1.5;color:#57606a;background-color:#fff8c5;border:1px solid #d4a72c;border-radius:6px;">%s</p>
</td></tr>
<tr><td style="padding:16px 28px 24px;background-color:#f6f8fa;border-top:1px solid #d0d7de;">
<p style="margin:0;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Helvetica,Arial,sans-serif;font-size:12px;line-height:1.5;color:#8b949e;">%s</p>
</td></tr>
</table>
</td></tr>
</table>
</body>
</html>`,
		esc(tpl.Title),
		esc(tpl.Intro),
		codeEsc,
		esc(tpl.ExpiryNote),
		esc(tpl.SecurityNote),
		esc(tpl.IgnoreNote),
		esc(tpl.Footer),
	)
}
