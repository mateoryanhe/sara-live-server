package call

import (
	"strconv"

	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	callentity "xr-game-server/entity/call"
	liveentity "xr-game-server/entity/live"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/module/liveroom"
)

type oneToOneRoomCallContext struct {
	audienceId     uint64
	orderParams    string
	pricePerMinute float64
}

func callSourceUsesLiveRoom(source uint8) bool {
	return source == callentity.CallOrderSourceLiveRoom ||
		source == callentity.CallOrderSourceOneToOneRoom
}

// resolveRoomCallPartyIdsByTypes 根据双方用户类型识别主播方和普通用户方。
// targetId 只是接听方，不保证一定是主播。
func resolveRoomCallPartyIdsByTypes(callerId uint64, callerType uint8, targetId uint64, targetType uint8) (anchorId, audienceId uint64, ok bool) {
	callerIsAnchor := userentity.UserTypeIsAnchor(callerType)
	targetIsAnchor := userentity.UserTypeIsAnchor(targetType)
	if callerIsAnchor == targetIsAnchor {
		return 0, 0, false
	}
	if callerIsAnchor {
		return callerId, targetId, true
	}
	return targetId, callerId, true
}

func resolveRoomCallPartyIds(callerId, targetId uint64) (anchorId, audienceId uint64, ok bool) {
	caller := userinfodao.GetUserInfoByUserId(callerId)
	target := userinfodao.GetUserInfoByUserId(targetId)
	if caller == nil || target == nil {
		return 0, 0, false
	}
	return resolveRoomCallPartyIdsByTypes(callerId, caller.UserType, targetId, target.UserType)
}

// resolveOneToOneRoomCallContext 根据双方用户类型找到主播房间，解析1v1房间呼叫配置。
func resolveOneToOneRoomCallContext(callerId, targetId uint64) (*oneToOneRoomCallContext, error) {
	anchorId, audienceId, ok := resolveRoomCallPartyIds(callerId, targetId)
	if !ok {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	room := liveroomdao.GetRoomById(anchorId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	cfg := liveroomdao.GetLiveRoomCfg(room.ID)
	if cfg == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	if cfg.Category != liveentity.LiveRoomCategoryOneToOne {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	if !liveroom.CanInitiateLiveRoomCall(room, cfg, audienceId) {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}

	return &oneToOneRoomCallContext{
		audienceId:     audienceId,
		orderParams:    strconv.FormatUint(room.LiveRecordId, 10),
		pricePerMinute: cfg.Billing,
	}, nil
}
