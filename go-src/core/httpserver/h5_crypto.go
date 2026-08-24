package httpserver

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"xr-game-server/core/xrjson"
	"xr-game-server/core/xrlog"
	"xr-game-server/errercode"
)

const (
	// H5ClientHeader 标识 H5 端请求/响应,值为 1 或 true 时启用 body 加解密
	H5ClientHeader = "X-H5-Client"
)

const (
	h5CryptoEnabledCtxKey = "httpserver.h5CryptoEnabled"
	h5DecryptMsCtxKey     = "httpserver.h5DecryptMs"
)

var h5DeploySecretProvider func() string

// SetH5DeploySecretProvider 注入 H5 部署密钥读取函数(由 h5livedeploy 模块初始化时注册)
func SetH5DeploySecretProvider(fn func() string) {
	h5DeploySecretProvider = fn
}

func isH5ClientRequest(r *ghttp.Request) bool {
	if r == nil {
		return false
	}
	v := strings.TrimSpace(r.GetHeader(H5ClientHeader))
	return v == "1" || strings.EqualFold(v, "true")
}

func markH5CryptoRequest(r *ghttp.Request) {
	if r == nil {
		return
	}
	r.SetCtxVar(h5CryptoEnabledCtxKey, true)
}

func isH5CryptoRequest(r *ghttp.Request) bool {
	if r == nil {
		return false
	}
	return r.GetCtxVar(h5CryptoEnabledCtxKey).Bool()
}

func stashH5DecryptMs(r *ghttp.Request, ms int64) {
	if r == nil {
		return
	}
	r.SetCtxVar(h5DecryptMsCtxKey, ms)
}

func h5DecryptMsFromRequest(r *ghttp.Request) int64 {
	if r == nil {
		return -1
	}
	return r.GetCtxVar(h5DecryptMsCtxKey).Int64()
}

func getH5DeploySecret() string {
	if h5DeploySecretProvider == nil {
		return ""
	}
	return strings.TrimSpace(h5DeploySecretProvider())
}

func deriveH5CryptoKey(secret string) []byte {
	sum := sha256.Sum256([]byte(secret))
	return sum[:]
}

// encryptH5Payload AES-256-GCM 加密,输出 base64(nonce+ciphertext)
func encryptH5Payload(secret string, plain []byte) ([]byte, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return nil, errors.New("h5 deploy secret is empty")
	}
	if len(plain) == 0 {
		return nil, nil
	}
	block, err := aes.NewCipher(deriveH5CryptoKey(secret))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, plain, nil)
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return []byte(encoded), nil
}

// decryptH5Payload 解密 base64(nonce+ciphertext) 请求体
func decryptH5Payload(secret string, encoded []byte) ([]byte, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return nil, errors.New("h5 deploy secret is empty")
	}
	encoded = bytesTrimSpace(encoded)
	if len(encoded) == 0 {
		return nil, nil
	}
	raw, err := base64.StdEncoding.DecodeString(string(encoded))
	if err != nil {
		return nil, fmt.Errorf("invalid h5 payload base64: %w", err)
	}
	block, err := aes.NewCipher(deriveH5CryptoKey(secret))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(raw) < nonceSize {
		return nil, errors.New("invalid h5 payload length")
	}
	nonce, ciphertext := raw[:nonceSize], raw[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func bytesTrimSpace(b []byte) []byte {
	return []byte(strings.TrimSpace(string(b)))
}

func writeH5ResponseBody(r *ghttp.Request, plain []byte, writeStart *gtime.Time) {
	start := gtime.Now()
	if writeStart != nil {
		start = writeStart
	}
	secret := getH5DeploySecret()
	encryptStart := gtime.Now()
	encrypted, err := encryptH5Payload(secret, plain)
	encryptMs := elapsedMs(encryptStart)
	if err != nil {
		xrlog.DetailLog.Errorf(r.Context(),
			"H5响应加密失败,reqId=%v,authId=%v,decryptMs=%vms,encryptMs=%vms,url=%v,err=%v",
			r.GetHeader(ReqId, ""),
			authIdFromRequest(r),
			h5DecryptMsFromRequest(r),
			encryptMs,
			r.RequestURI,
			err,
		)
		fail := CreateFail(errercode.SysError)
		fail.Message = err.Error()
		plain = xrjson.MustMarshal(fail)
		encrypted = plain
		encryptMs = 0
	}
	r.Response.Header().Set(H5ClientHeader, "1")
	r.Response.Header().Set("Content-Type", contentTypeJson)
	r.Response.Write(encrypted)
	stashAPIResponseBufferWrittenAt(r)
	logAPIRequestResponseWrite(r, elapsedMs(start), len(encrypted), string(plain))
	logH5Crypto(r, h5DecryptMsFromRequest(r), encryptMs, len(plain), len(encrypted))
}

func logH5Crypto(r *ghttp.Request, decryptMs, encryptMs int64, plainBytes, encryptedBytes int) {
	if r == nil || !isH5CryptoRequest(r) {
		return
	}
	xrlog.DetailLog.Infof(r.Context(),
		"H5加解密完成,reqId=%v,authId=%v,decryptMs=%vms,encryptMs=%vms,plainBytes=%v,encryptedBytes=%v,url=%v",
		r.GetHeader(ReqId, ""),
		authIdFromRequest(r),
		decryptMs,
		encryptMs,
		plainBytes,
		encryptedBytes,
		r.RequestURI,
	)
}
