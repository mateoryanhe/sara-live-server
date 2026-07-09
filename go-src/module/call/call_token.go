package call

import (
	"context"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/calldto"
	"xr-game-server/errercode"
	"xr-game-server/module/agora"
)

// GetCallToken App端上报 channelName 获取通话频道 Token
func GetCallToken(ctx context.Context, req *calldto.CallTokenReq) (*calldto.CallTokenRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}

	token, tokenExpireAt, err := agora.ResolveCallToken(userId, req.ChannelName)
	if err != nil {
		return nil, err
	}

	appId := ""
	if agoraCfg, err := agora.GetAppId(ctx, nil); err == nil && agoraCfg != nil {
		appId = agoraCfg.AppId
	}

	return &calldto.CallTokenRes{
		Token:       token,
		AppId:       appId,
		ChannelName: req.ChannelName,
		UserAccount: strconv.FormatUint(userId, 10),
		ExpireAt:    tokenExpireAt,
	}, nil
}
