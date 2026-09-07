package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/authdto"
)

// SendEmailCode App 发邮箱验证码:服务端生成 6 位码并走 Cloudflare 发信(含冷却/日限与验证码缓存).
func SendEmailCode(ctx context.Context, req *authdto.SendEmailCodeReq) (*authdto.SendEmailCodeRes, error) {
	lang := req.Lang
	if lang == "" {
		if r := g.RequestFromCtx(ctx); r != nil {
			lang = r.GetHeader(httpserver.AcceptLanguageHeader)
		}
	}
	code, err := generateEmailVerifyCode()
	if err != nil {
		return nil, err
	}
	if err := SendEmailVerifyCode(ctx, req.Email, code, lang); err != nil {
		return nil, err
	}
	return &authdto.SendEmailCodeRes{Success: true}, nil
}

func generateEmailVerifyCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
