package shareappdownload

import (
	"net/url"
	"strings"

	"xr-game-server/core/cfg"
	"xr-game-server/dto/officialsitedeploydto"
)

// ApkShortPath 分享短链路径: GET 后 302 到官网 assets 下当前最大 APK.
const ApkShortPath = "/dl"

// InviteCodeLinePrefix 分享内容中邀请码行的标题前缀.
const InviteCodeLinePrefix = "Invite Code: "

// PreviewInviteCodePlaceholder CMS 预览用的示例邀请码.
const PreviewInviteCodePlaceholder = "YOUR_INVITE_CODE"

// GetApkShortUrl 官网域名下的短链(有 APK 时才返回,否则空).
func GetApkShortUrl() string {
	if GetDownloadUrl() == "" {
		return ""
	}
	domain := strings.TrimSpace(cfg.GetStaticSiteDomain(officialsitedeploydto.OfficialSiteStaticPrefix))
	if domain == "" {
		return ""
	}
	return (&url.URL{
		Scheme: "https",
		Host:   domain,
		Path:   ApkShortPath,
	}).String()
}

// FormatInviteCodeLine 邀请码行(带标题,方便识别).
func FormatInviteCodeLine(inviteCode string) string {
	inviteCode = strings.TrimSpace(inviteCode)
	if inviteCode == "" {
		return ""
	}
	return InviteCodeLinePrefix + inviteCode
}

// BuildAppShareCode App getUserInfo.shareCode:
// 第1段分享文案 / 第2段「Invite Code: 邀请码」 / 第3段短链(点开重定向到最新最大 APK).
func BuildAppShareCode(inviteCode string) string {
	inviteCode = strings.TrimSpace(inviteCode)
	text := strings.TrimSpace(GetShareText())
	// 旧占位符不再参与拼接,避免出现在正文里
	text = strings.ReplaceAll(text, "{inviteCode}", "")
	text = strings.ReplaceAll(text, "{downloadUrl}", "")
	text = strings.TrimSpace(text)

	parts := make([]string, 0, 3)
	if text != "" {
		parts = append(parts, strings.TrimRight(text, "\r\n"))
	}
	if line := FormatInviteCodeLine(inviteCode); line != "" {
		parts = append(parts, line)
	}
	if shortURL := GetApkShortUrl(); shortURL != "" {
		parts = append(parts, shortURL)
	}
	return strings.Join(parts, "\n")
}

// BuildPreviewShareCode CMS 展示用:与 App 同规则,邀请码用占位符.
func BuildPreviewShareCode() string {
	return BuildAppShareCode(PreviewInviteCodePlaceholder)
}
