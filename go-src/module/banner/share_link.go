package banner

import (
	"net/url"
	"strings"
)

// appendShareParamsToLink Banner 跳转链接追加 token(当前用户 App Token,便于 H5 鉴权查数据).
func appendShareParamsToLink(rawLink, appToken string) string {
	link := strings.TrimSpace(rawLink)
	appToken = strings.TrimSpace(appToken)
	if link == "" || appToken == "" {
		return link
	}

	u, err := url.Parse(link)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return appendShareParamsRaw(link, appToken)
	}
	q := u.Query()
	q.Set("token", appToken)
	u.RawQuery = q.Encode()
	return u.String()
}

func appendShareParamsRaw(link, appToken string) string {
	sep := "?"
	if strings.Contains(link, "?") {
		sep = "&"
	}
	return link + sep + "token=" + url.QueryEscape(appToken)
}
