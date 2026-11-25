package shareappdownload

import (
	"context"
	"net/url"
	"strings"

	"xr-game-server/core/cfg"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/officialsitedeploydto"
)

const (
	// SofieDefaultShortPath 默认配置(无包名)短链路径.
	SofieDefaultShortPath = "/dl/sofie"
	// SofiePackageShortPathPattern 按包名短链: /dl/p/{packageName}
	SofiePackageShortPathPattern = "/dl/p/{packageName}"
)

func sofieShareBaseDomain() string {
	d := strings.TrimSpace(cfg.GetStaticSiteDomain(officialsitedeploydto.SofieOfficialSiteStaticPrefix))
	if d != "" {
		return d
	}
	return strings.TrimSpace(cfg.GetStaticSiteDomain(officialsitedeploydto.OfficialSiteStaticPrefix))
}

// BuildSofieShortUrl 按包名生成分享短链; 默认行使用 /dl/sofie.
func BuildSofieShortUrl(packageName string) string {
	pkg := strings.TrimSpace(packageName)
	if ResolveSofieGooglePlayUrl(pkg) == "" {
		return ""
	}
	domain := sofieShareBaseDomain()
	if domain == "" {
		return ""
	}
	path := SofieDefaultShortPath
	if pkg != "" {
		path = "/dl/p/" + url.PathEscape(pkg)
	}
	return (&url.URL{
		Scheme: "https",
		Host:   domain,
		Path:   path,
	}).String()
}

func buildShareCodeFromSofieSnapshot(snap *sofieShareRowSnapshot, inviteCode string) string {
	if snap == nil {
		return ""
	}
	inviteCode = strings.TrimSpace(inviteCode)
	text := strings.TrimSpace(snap.ShareText)
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
	if shortURL := BuildSofieShortUrl(snap.PackageName); shortURL != "" {
		parts = append(parts, shortURL)
	}
	return strings.Join(parts, "\n")
}

// BuildAppShareCodeForApp 按请求(Sofie 路径 + X-App-Package)下发分享文案; 否则回落 Sara.
func BuildAppShareCodeForApp(ctx context.Context, inviteCode string) string {
	if httpserver.IsSofieAppRequest(ctx) {
		pkg := httpserver.GetPackageNameFromContext(ctx)
		snap := ResolveSofieShareSnapshot(pkg)
		if snap != nil && (snap.GooglePlayUrl != "" || snap.ShareText != "" || BuildSofieShortUrl(snap.PackageName) != "") {
			if built := buildShareCodeFromSofieSnapshot(snap, inviteCode); built != "" {
				return built
			}
		}
		// Sofie 未配置时仍尝试默认行
		if def := ResolveSofieShareSnapshot(""); def != nil {
			if built := buildShareCodeFromSofieSnapshot(def, inviteCode); built != "" {
				return built
			}
		}
	}
	return BuildAppShareCode(inviteCode)
}

// BuildSofiePreviewShareCode CMS 预览.
func BuildSofiePreviewShareCode(packageName string) string {
	snap := ResolveSofieShareSnapshot(strings.TrimSpace(packageName))
	if snap == nil {
		snap = ResolveSofieShareSnapshot("")
	}
	return buildShareCodeFromSofieSnapshot(snap, PreviewInviteCodePlaceholder)
}

// ParseSofiePackageFromShortPath 从 /dl/p/{packageName} 解析包名.
func ParseSofiePackageFromShortPath(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	decoded, err := url.PathUnescape(raw)
	if err != nil {
		decoded = raw
	}
	return strings.TrimSpace(decoded)
}