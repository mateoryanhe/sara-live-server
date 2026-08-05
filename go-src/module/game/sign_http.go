package game

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"strings"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/module/game/sign"
)

// BuildSign 生成 X-Sign.
func BuildSign(params map[string]string, secretKey string) string {
	return sign.Build(params, secretKey)
}

// VerifySign 校验签名.
func VerifySign(params map[string]string, signValue, secretKey string) bool {
	return sign.Verify(params, signValue, secretKey)
}

// VerifySignWithPlatformCfg 使用 CMS 平台 SecretKey 校验签名.
func VerifySignWithPlatformCfg(params map[string]string, signValue string) bool {
	cfg := cfgdao.GetGamePlatformCfgFromMemory()
	if cfg == nil {
		return false
	}
	return sign.Verify(params, signValue, cfg.SecretKey)
}

// MergeSignParams 合并 query、body 与 Header 参与签名的字段.
func MergeSignParams(query map[string]string, body map[string]string, operatorToken, timestamp string) map[string]string {
	return sign.MergeParams(query, body, operatorToken, timestamp)
}

// CollectSignParamsFromRequest 从 HTTP 请求收集参与签名的参数(不含 sign).
func CollectSignParamsFromRequest(r *ghttp.Request) (map[string]string, string) {
	if r == nil {
		return nil, ""
	}
	query := sign.ParseQuery(r.URL.Query())
	body := sign.ParseJSONBody(r.GetBody())
	operatorToken := r.Header.Get(sign.HeaderOperatorToken)
	timestamp := r.Header.Get(sign.HeaderTimestamp)
	signValue := strings.TrimSpace(query[sign.FieldSign])
	if signValue == "" {
		signValue = strings.TrimSpace(body[sign.FieldSign])
	}
	delete(query, sign.FieldSign)
	delete(body, sign.FieldSign)
	return sign.MergeParams(query, body, operatorToken, timestamp), signValue
}

// VerifyRequestSign 校验 HTTP 回调请求签名.
func VerifyRequestSign(r *ghttp.Request, secretKey string) bool {
	params, signValue := CollectSignParamsFromRequest(r)
	return sign.Verify(params, signValue, secretKey)
}

// VerifyRequestSignWithPlatformCfg 使用平台配置 SecretKey 校验 HTTP 回调请求签名.
func VerifyRequestSignWithPlatformCfg(r *ghttp.Request) bool {
	params, signValue := CollectSignParamsFromRequest(r)
	return VerifySignWithPlatformCfg(params, signValue)
}
