package agora

import (
	"context"
	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/agoradto"
	"xr-game-server/errercode"
	"xr-game-server/module/livecfg"
)

// BuildLiveRoomToken 为指定用户生成进直播间声网 Token
func BuildLiveRoomToken(userId, roomId uint64) (token string, expireAt int64, err error) {
	if liveroomdao.GetRoomById(roomId) == nil {
		return "", 0, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	agoraCfg := getAgoraCfgCache()
	if err := validateAgoraCfg(agoraCfg); err != nil {
		return "", 0, err
	}

	channelName := buildChannelName(roomId)
	userAccount := buildUserAccount(userId)
	role := resolveRole(userId, roomId)
	expireSeconds := agoraCfg.TokenExpireSeconds

	token, err = rtctokenbuilder.BuildTokenWithUserAccount(
		agoraCfg.AppId,
		agoraCfg.AppCertificate,
		channelName,
		userAccount,
		role,
		expireSeconds,
		expireSeconds,
	)
	if err != nil {
		return "", 0, err
	}

	return token, int64(expireSeconds), nil
}

// GetLiveRoomToken App端上报房间ID,返回进直播间声网Token
func GetLiveRoomToken(ctx context.Context, req *agoradto.GetLiveRoomTokenReq) (*agoradto.GetLiveRoomTokenRes, error) {
	userId := httpserver.GetAuthId(ctx)
	token, expireAt, err := BuildLiveRoomToken(userId, req.RoomId)
	if err != nil {
		return nil, err
	}

	agoraCfg := getAgoraCfgCache()
	return &agoradto.GetLiveRoomTokenRes{
		Token:       token,
		AppId:       agoraCfg.AppId,
		ChannelName: buildChannelName(req.RoomId),
		UserAccount: buildUserAccount(userId),
		ExpireAt:    expireAt,
	}, nil
}

// GetAppId App端获取声网AppId
func GetAppId(_ context.Context, _ *agoradto.GetAppIdReq) (*agoradto.GetAppIdRes, error) {
	agoraCfg := getAgoraCfgCache()
	if agoraCfg == nil || agoraCfg.AppId == "" {
		return nil, errercode.CreateCode(errercode.AgoraCfgInvalid)
	}
	return &agoradto.GetAppIdRes{
		AppId:                       agoraCfg.AppId,
		PrivateRoomFreeWatchSeconds: livecfg.GetPrivateRoomFreeWatchSeconds(),
		PaidDanmakuPrice:            livecfg.GetPaidDanmakuPrice(),
	}, nil
}

func buildChannelName(roomId uint64) string {
	return strconv.FormatUint(roomId, 10)
}

func buildUserAccount(userId uint64) string {
	return strconv.FormatUint(userId, 10)
}

func resolveRole(userId, roomId uint64) rtctokenbuilder.Role {
	if userId == roomId {
		return rtctokenbuilder.RolePublisher
	}
	return rtctokenbuilder.RoleSubscriber
}
