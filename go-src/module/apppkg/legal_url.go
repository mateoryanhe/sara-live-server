package apppkg

// ResolvePrivacyPolicyUrl 隐私政策 URL：已取消包级覆盖，直接使用全局配置 fallback。
func ResolvePrivacyPolicyUrl(_ string, fallback string) string {
	return fallback
}

// ResolveTermsOfServiceUrl 用户协议 URL：已取消包级覆盖，直接使用全局配置 fallback。
func ResolveTermsOfServiceUrl(_ string, fallback string) string {
	return fallback
}
