package httpserver

import (
	"bytes"
	"io"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"xr-game-server/core/xrlog"
	"xr-game-server/errercode"
)

// MiddlewareH5Crypto H5 端 body 加解密中间件: 请求头 X-H5-Client=1/true 时生效
func MiddlewareH5Crypto(r *ghttp.Request) {
	if !isH5ClientRequest(r) {
		r.Middleware.Next()
		return
	}
	markH5CryptoRequest(r)

	secret := getH5DeploySecret()
	if secret == "" {
		xrlog.DetailLog.Errorf(r.Context(), "H5加解密失败,部署密钥未配置,url=%v", r.RequestURI)
		WriteFailJson(r, int(errercode.InvalidParam))
		return
	}

	if r.Request != nil && r.Request.Body != nil && r.ContentLength != 0 {
		decryptStart := gtime.Now()
		encryptedBody, err := io.ReadAll(r.Request.Body)
		if err != nil {
			stashH5DecryptMs(r, elapsedMs(decryptStart))
			xrlog.DetailLog.Errorf(r.Context(), "H5请求体读取失败,url=%v,err=%v", r.RequestURI, err)
			WriteFailJson(r, int(errercode.InvalidParam))
			return
		}
		plainBody, err := decryptH5Payload(secret, encryptedBody)
		decryptMs := elapsedMs(decryptStart)
		stashH5DecryptMs(r, decryptMs)
		if err != nil {
			xrlog.DetailLog.Errorf(r.Context(),
				"H5请求体解密失败,reqId=%v,decryptMs=%vms,url=%v,err=%v",
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
