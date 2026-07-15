package agora

import (
	"context"
	"strconv"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/agoradto"
	"xr-game-server/errercode"
	"xr-game-server/module/livecfg"
)

func validateAgoraCfg(cfg *agoraCfgSnapshot) error {
	if cfg == nil || cfg.AppId == "" || cfg.AppCertificate == "" {
		return errercode.CreateCode(errercode.AgoraCfgInvalid)
	}
	return nil
}

// BuildChannelToken 为指定频道生成声网 Token(频道名与角色由业务方决定)
func BuildChannelToken(userId uint64, channelName string, role uint8) (token string, expireAt int64, err error) {
	rtcRole, err := toRTCRole(role)
	if err != nil {
		return "", 0, err
	}

	agoraCfg := getAgoraCfgCache()
	if err := validateAgoraCfg(agoraCfg); err != nil {
		return "", 0, err
	}

	userAccount := buildUserAccount(userId)
	expireSeconds := getChannelTokenExpireSeconds()
	token, err = rtctokenbuilder.BuildTokenWithUserAccount(
		agoraCfg.AppId,
		agoraCfg.AppCertificate,
		channelName,
		userAccount,
		rtcRole,
		expireSeconds,
		expireSeconds,
	)
	if err != nil {
		return "", 0, err
	}
	return token, int64(expireSeconds), nil
}

func toRTCRole(role uint8) (rtctokenbuilder.Role, error) {
	switch role {
	case agoradto.RTCRolePublisher:
		return rtctokenbuilder.RolePublisher, nil
	case agoradto.RTCRoleSubscriber:
		return rtctokenbuilder.RoleSubscriber, nil
	default:
		return 0, errercode.CreateCode(errercode.InvalidParam)
	}
}

// GetLiveRoomToken App端上报频道名与角色,返回最新声网Token
func GetLiveRoomToken(ctx context.Context, req *agoradto.GetLiveRoomTokenReq) (*agoradto.GetLiveRoomTokenRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	token, expireAt, err := ResolveChannelToken(userId, req.ChannelName, req.Role)
	if err != nil {
		return nil, err
	}

	agoraCfg := getAgoraCfgCache()
	return &agoradto.GetLiveRoomTokenRes{
		Token:       token,
		AppId:       agoraCfg.AppId,
		ChannelName: req.ChannelName,
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

func buildUserAccount(userId uint64) string {
	return strconv.FormatUint(userId, 10)
}
