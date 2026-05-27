package agora

import (
	"context"
	"math"

	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/agoradto"
	"xr-game-server/errercode"
)

// CheckUserOnline 根据房间ID与用户ID,调用声网官方 REST API 查询用户是否在频道内
func CheckUserOnline(ctx context.Context, req *agoradto.CheckUserOnlineReq) (*agoradto.CheckUserOnlineRes, error) {
	if liveroomdao.GetRoomById(req.RoomId) == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	cfg := getAgoraCfgCache()
	if err := validateAgoraRestCfg(cfg); err != nil {
		return nil, err
	}

	channelName := buildChannelName(req.RoomId)
	userAccount := buildUserAccount(req.UserId)

	online, role, joinAt, err := queryUserOnlineInChannel(ctx, cfg, channelName, req.UserId, userAccount)
	if err != nil {
		return nil, err
	}

	return &agoradto.CheckUserOnlineRes{
		Online:      online,
		ChannelName: channelName,
		UserAccount: userAccount,
		Role:        role,
		JoinAt:      joinAt,
	}, nil
}

func validateAgoraCfg(cfg *agoraCfgSnapshot) error {
	if cfg == nil || cfg.AppId == "" || cfg.AppCertificate == "" {
		return errercode.CreateCode(errercode.AgoraCfgInvalid)
	}
	return nil
}

func validateAgoraRestCfg(cfg *agoraCfgSnapshot) error {
	if err := validateAgoraCfg(cfg); err != nil {
		return err
	}
	if cfg.RestCustomerId == "" || cfg.RestCustomerSecret == "" {
		return errercode.CreateCode(errercode.AgoraCfgInvalid)
	}
	return nil
}

func queryUserOnlineInChannel(
	ctx context.Context,
	cfg *agoraCfgSnapshot,
	channelName string,
	userId uint64,
	userAccount string,
) (online bool, role int, joinAt int64, err error) {
	// 优先走频道用户列表(兼容 UserAccount 进房)
	listData, err := queryChannelUserList(ctx, cfg, channelName)
	if err != nil {
		return false, 0, 0, err
	}
	if userAccountInChannelLists(userAccount, listData) {
		return true, 0, 0, nil
	}

	// 兜底: 若 userId 在 uint32 范围内,再查官方 user/property 接口
	if userId > 0 && userId <= math.MaxUint32 {
		status, queryErr := queryUserStatus(ctx, cfg, channelName, userId)
		if queryErr != nil {
			return false, 0, 0, queryErr
		}
		if status != nil && status.InChannel {
			return true, status.Role, status.Join, nil
		}
	}
	return false, 0, 0, nil
}
