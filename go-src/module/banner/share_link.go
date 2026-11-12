package banner

import (
	"net/url"
	"strings"
)

// appendShareParamsToLink Banner 跳转链接追加 token(登录同款 userId.token,便于 H5 鉴权).
func appendShareParamsToLink(rawLink, clientToken string) string {
	link := strings.TrimSpace(rawLink)
	clientToken = strings.TrimSpace(clientToken)
	if link == "" || clientToken == "" {
		return link
	}

	u, err := url.Parse(link)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return appendShareParamsRaw(link, clientToken)
	}
	q := u.Query()
	q.Set("token", clientToken)
	u.RawQuery = q.Encode()
	return u.String()
}

func appendShareParamsRaw(link, clientToken string) string {
	sep := "?"
	if strings.Contains(link, "?") {
		sep = "&"
	}
	return link + sep + "token=" + url.QueryEscape(clientToken)
}
