package apppkg

import "strings"

// ResolvePrivacyPolicyUrl 按包名获取隐私政策 URL,包级未配置时使用 fallback
func ResolvePrivacyPolicyUrl(packageName, fallback string) string {
	pkg := GetAppPkgFromMemoryByPackageName(packageName)
	if pkg != nil {
		if url := strings.TrimSpace(pkg.PrivacyPolicyUrl); url != "" {
			return url
		}
	}
	return fallback
}

// ResolveTermsOfServiceUrl 按包名获取用户服务协议 URL,包级未配置时使用 fallback
func ResolveTermsOfServiceUrl(packageName, fallback string) string {
	pkg := GetAppPkgFromMemoryByPackageName(packageName)
	if pkg != nil {
		if url := strings.TrimSpace(pkg.TermsOfServiceUrl); url != "" {
			return url
		}
	}
	return fallback
}
