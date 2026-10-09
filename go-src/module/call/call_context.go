package call

import (
	"strconv"

	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	callentity "xr-game-server/entity/call"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/module/liveroom"
)

type oneToOneRoomCallContext struct {
	anchorId       uint64
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
	caller := userinfodao.GetUserInfoByUserId(callerId)
	target := userinfodao.GetUserInfoByUserId(targetId)
	if caller == nil || target == nil {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	callerIsAnchor := userentity.UserTypeIsAnchor(caller.UserType)
	targetIsAnchor := userentity.UserTypeIsAnchor(target.UserType)
	if callerIsAnchor && targetIsAnchor {
		return nil, errercode.CreateCode(errercode.OneToOneRoomCallBothAnchors)
	}
	anchorId, audienceId, ok := resolveRoomCallPartyIdsByTypes(callerId, caller.UserType, targetId, target.UserType)
	if !ok {
		return nil, errercode.CreateCode(errercode.NoPermission)
	}
	if !liveroom.IsOneToOneOnShelf(anchorId) {
		return nil, errercode.CreateCode(errercode.LiveRoomNotOneToOne)
	}
	oneToOne := liveroomdao.GetOneToOneRoom(anchorId)
	if oneToOne == nil || !oneToOne.IsOnShelf() {
		return nil, errercode.CreateCode(errercode.LiveRoomNotOneToOne)
	}
	// 1v1(source=3)不要求 live_rooms 存在，也不校验 privateInviteType；有开播场次时仍写入 params 供收益关联。
	orderParams := "0"
	if room := liveroomdao.ResolveRoom(anchorId); room != nil && room.LiveRecordId > 0 {
		orderParams = strconv.FormatUint(room.LiveRecordId, 10)
	}

	return &oneToOneRoomCallContext{
		anchorId:       anchorId,
		audienceId:     audienceId,
		orderParams:    orderParams,
		pricePerMinute: oneToOne.Billing,
	}, nil
}
