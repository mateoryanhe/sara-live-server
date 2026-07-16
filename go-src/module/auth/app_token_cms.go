package auth

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrtoken"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dto/apptokendto"
	"xr-game-server/entity"
)

const defaultAppTokenExpireDays = 30

func GetAppToken(ctx context.Context, req *apptokendto.GetAppTokenReq) (*httpserver.CMSQueryResp, error) {
	_ = ctx
	userId := gconv.Uint64(req.UserId)
	total, list := accountdao.QueryAppTokens(userId, req.PageIndex, req.PageSize)
	data := make([]*apptokendto.AppTokenDto, 0, len(list))
	for _, item := range list {
		data = append(data, apptokendto.NewAppTokenDto(item))
	}
	return httpserver.NewCMSQueryResp(total, data), nil
}

func SaveAppToken(ctx context.Context, req *apptokendto.SaveAppTokenReq) (bool, error) {
	_ = ctx
	token := strings.TrimSpace(req.Token)
	existing := accountdao.GetAppTokenByUserId(req.Id)
	if token == "" {
		if existing != nil && existing.Token != "" {
			token = existing.Token
		} else {
			token = guid.S()
		}
	}

	expireAt := time.Now().Add(defaultAppTokenExpireDays * 24 * time.Hour)
	if req.ExpireAt != nil {
		expireAt = *req.ExpireAt
	} else if existing != nil && !existing.ExpireAt.IsZero() {
		expireAt = existing.ExpireAt
	}

	xrtoken.InitAppToken(req.Id, token, expireAt)
	entity.NewAppToken(req.Id, token, expireAt)
	return true, nil
}
