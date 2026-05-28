package auth

import (
	"context"

	"github.com/gogf/gf/v2/util/gconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrtoken"
	"xr-game-server/dto/apptokendto"
	"xr-game-server/entity"
)

func GetAppToken(ctx context.Context, req *apptokendto.GetAppTokenReq) (*httpserver.CMSQueryResp, error) {
	userId := gconv.Uint64(req.UserId)
	total, list := xrtoken.QueryAppTokens(userId, req.PageIndex, req.PageSize)
	data := make([]*apptokendto.AppTokenDto, 0, len(list))
	for _, item := range list {
		data = append(data, apptokendto.NewAppTokenDtoFromCache(item.UserId, item.Token, item.ExpireAt))
	}
	return httpserver.NewCMSQueryResp(total, data), nil
}

func SaveAppToken(ctx context.Context, req *apptokendto.SaveAppTokenReq) (bool, error) {

	xrtoken.InitAppToken(req.Id, req.Token, *req.ExpireAt)
	entity.NewAppToken(req.Id, req.Token, *req.ExpireAt)
	return true, nil
}
