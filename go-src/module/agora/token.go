package agora

import (
	"context"
	"strconv"
	"time"

	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/agoradto"
	"xr-game-server/errercode"
	"xr-game-server/module/livecfg"
)

const subscriberTokenCacheTTL = 10 * time.Minute

var subscriberTokenCache = gcache.New()

type channelTokenResult struct {
	Token    string
	ExpireAt int64
}

func validateAgoraCfg(cfg *agoraCfgSnapshot) error {
	if cfg == nil || cfg.AppId == "" || cfg.AppCertificate == "" {
		return errercode.CreateCode(errercode.AgoraCfgInvalid)
	}
	return nil
}

// BuildChannelToken 为指定频道生成声网 Token
// RoleSubscriber: 生成通配 Token(空频道名), gcache 缓存 10 分钟
// RolePublisher: 每次生成新 Token, 不缓存
func BuildChannelToken(userId uint64, channelName string, role uint8) (token string, expireAt int64, err error) {
	rtcRole, err := toRTCRole(role)
	if err != nil {
		return "", 0, err
	}

	agoraCfg := getAgoraCfgCache()
	if err := validateAgoraCfg(agoraCfg); err != nil {
		return "", 0, err
	}

	if role == agoradto.RTCRoleSubscriber {
		return getSubscriberWildcardToken(userId, rtcRole)
	}

	if channelName == "" {
		return "", 0, errercode.CreateCode(errercode.InvalidParam)
	}
	result, err := buildChannelToken(userId, channelName, rtcRole)
	if err != nil {
		return "", 0, err
	}
	return result.Token, result.ExpireAt, nil
}

func getSubscriberWildcardToken(userId uint64, rtcRole rtctokenbuilder.Role) (token string, expireAt int64, err error) {
	ctx := gctx.New()
	v, err := subscriberTokenCache.GetOrSetFuncLock(ctx, userId, func(ctx context.Context) (interface{}, error) {
		return buildChannelToken(userId, "", rtcRole)
	}, subscriberTokenCacheTTL)
	if err != nil {
		return "", 0, err
	}
	result, ok := v.Val().(*channelTokenResult)
	if !ok || result == nil {
		return "", 0, errercode.CreateCode(errercode.SysError)
	}
	return result.Token, result.ExpireAt, nil
}

func buildChannelToken(userId uint64, channelName string, rtcRole rtctokenbuilder.Role) (*channelTokenResult, error) {
	agoraCfg := getAgoraCfgCache()
	expireSeconds := getChannelTokenExpireSeconds()
	token, err := rtctokenbuilder.BuildTokenWithUserAccount(
		agoraCfg.AppId,
		agoraCfg.AppCertificate,
		channelName,
		buildUserAccount(userId),
		rtcRole,
		expireSeconds,
		expireSeconds,
	)
	if err != nil {
		return nil, err
	}
	return &channelTokenResult{
		Token:    token,
		ExpireAt: time.Now().Add(time.Duration(expireSeconds) * time.Second).Unix(),
	}, nil
}

func clearSubscriberTokenCache() {
	_ = subscriberTokenCache.Clear(gctx.New())
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

// ResolveChannelToken 按频道名查询声网 Token(频道名与角色由业务方决定)
func ResolveChannelToken(userId uint64, channelName string, role uint8) (token string, expireAt int64, err error) {
	if channelName == "" {
		return "", 0, errercode.CreateCode(errercode.InvalidParam)
	}
	return BuildChannelToken(userId, channelName, role)
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
