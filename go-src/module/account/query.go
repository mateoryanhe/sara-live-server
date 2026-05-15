package account

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dto/accountdto"
)

func QueryUserInfo(ctx context.Context, req *accountdto.QueryUserInfoReq) (res *httpserver.CMSQueryResp, err error) {
	total, data := accountdao.GetUserInfo(req)
	return httpserver.NewCMSQueryResp(total, data), nil
}
