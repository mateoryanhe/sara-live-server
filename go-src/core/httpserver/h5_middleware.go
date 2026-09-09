package httpserver

import (
	"bytes"
	"io"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"xr-game-server/core/xrlog"
	"xr-game-server/errercode"
)

// MiddlewareH5Crypto 客户端 body 加解密中间件:
// - Header X-H5-Client=1/true → 使用 H5 直播部署密钥
// - Header X-Coin-Merchant-Client=1/true → 使用币商部署密钥
// 未带上述 Header 时明文透传。
func MiddlewareH5Crypto(r *ghttp.Request) {
	kind := ""
	switch {
	case isH5ClientRequest(r):
		kind = clientCryptoKindH5
	case isCoinMerchantClientRequest(r):
		kind = clientCryptoKindCoinMer
	default:
		r.Middleware.Next()
		return
	}
	markH5CryptoRequest(r, kind)

	secret := getClientCryptoSecret(r)
	if secret == "" {
		xrlog.DetailLog.Errorf(r.Context(), "客户端加解密失败,部署密钥未配置,kind=%v,url=%v", kind, r.RequestURI)
		WriteFailJson(r, int(errercode.InvalidParam))
		return
	}

	if r.Request != nil && r.Request.Body != nil && r.ContentLength != 0 {
		decryptStart := gtime.Now()
		encryptedBody, err := io.ReadAll(r.Request.Body)
		if err != nil {
			stashH5DecryptMs(r, elapsedMs(decryptStart))
			xrlog.DetailLog.Errorf(r.Context(), "客户端请求体读取失败,kind=%v,url=%v,err=%v", kind, r.RequestURI, err)
			WriteFailJson(r, int(errercode.InvalidParam))
			return
		}
		plainBody, err := decryptH5Payload(secret, encryptedBody)
		decryptMs := elapsedMs(decryptStart)
		stashH5DecryptMs(r, decryptMs)
		if err != nil {
			xrlog.DetailLog.Errorf(r.Context(),
				"客户端请求体解密失败,kind=%v,reqId=%v,decryptMs=%vms,url=%v,err=%v",
				kind,
				r.GetHeader(ReqId, ""),
				decryptMs,
				r.RequestURI,
				err,
			)
			WriteFailJson(r, int(errercode.H5PayloadDecodeFail))
			return
		}
		r.Request.Body = io.NopCloser(bytes.NewReader(plainBody))
		r.Request.ContentLength = int64(len(plainBody))
		if r.Request.ContentLength > 0 {
			r.Header.Set("Content-Type", contentTypeJson)
		}
	}

	r.Middleware.Next()
}
