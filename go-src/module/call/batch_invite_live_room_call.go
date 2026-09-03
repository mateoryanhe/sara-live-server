package call

import (
	"context"
	"strconv"
	"strings"

	"xr-game-server/constants/cmd"
	"xr-game-server/core/httpserver"
	"xr-game-server/core/push"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/calldto"
	callentity "xr-game-server/entity/call"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

const (
	batchInviteLiveRoomCallMax = 100
	liveRoomCallInviteMessage  = "主播邀请你通话"
)

// BatchInviteLiveRoomCall 主播批量邀请本房在线观众通话(仅推送,不建通话单;观众可再调 liveRoomCall 发起)
func BatchInviteLiveRoomCall(ctx context.Context, req *calldto.BatchInviteLiveRoomCallReq) (*calldto.BatchInviteLiveRoomCallRes, error) {
	anchorId := httpserver.GetAuthId(ctx)
	if anchorId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if len(req.UserIds) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if len(req.UserIds) > batchInviteLiveRoomCallMax {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	room := liveroomdao.GetRoomById(anchorId)
	if room == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}
	if room.LiveRecordId == 0 {
		return nil, errercode.CreateCode(errercode.LiveRoomNotLive)
	}
	cfg := liveroomdao.GetLiveRoomCfg(room.ID)
	if cfg == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomNotExist)
	}

	seen := make(map[uint64]struct{}, len(req.UserIds))
	targetIds := make([]uint64, 0, len(req.UserIds))
	skipped := make([]string, 0)
	for _, raw := range req.UserIds {
		idStr := strings.TrimSpace(raw)
		if idStr == "" {
			continue
		}
		userId, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil || userId == 0 {
			skipped = append(skipped, idStr)
			continue
		}
		if userId == anchorId {
			skipped = append(skipped, idStr)
			continue
		}
		if _, ok := seen[userId]; ok {
			continue
		}
		seen[userId] = struct{}{}
		targetIds = append(targetIds, userId)
	}

	item := buildLiveRoomCallInviteItem(anchorId, cfg)
	invited := 0
	for _, userId := range targetIds {
		onlineId := liveentity.BuildLiveRoomOnlineId(userId, anchorId)
		online := liveroomdao.GetOnlineById(onlineId, userId, anchorId)
		if online == nil || online.Status != liveentity.LiveRoomOnlineStatusOnline {
			skipped = append(skipped, strconv.FormatUint(userId, 10))
			continue
		}
		if userinfodao.GetUserInfoByUserId(userId) == nil {
			skipped = append(skipped, strconv.FormatUint(userId, 10))
			continue
		}
		push.Data(userId, cmd.LiveRoomCallInvite, item)
		invited++
	}

	return &calldto.BatchInviteLiveRoomCallRes{
		InvitedCount:   invited,
		SkippedUserIds: skipped,
	}, nil
}

func buildLiveRoomCallInviteItem(anchorId uint64, cfg *liveentity.LiveRoomCfg) *calldto.CallInvitePushItem {
	item := &calldto.CallInvitePushItem{
		RoomId:   strconv.FormatUint(anchorId, 10),
		AnchorId: strconv.FormatUint(anchorId, 10),
		CallType: callentity.CallOrderTypeVideo,
		Message:  liveRoomCallInviteMessage,
	}
	if cfg != nil {
		item.Ticket = cfg.Ticket
		item.Billing = cfg.Billing
	}
	if u := userinfodao.GetUserInfoByUserId(anchorId); u != nil {
		item.AnchorNickname = u.Nickname
		item.AnchorAvatar = upload.ResolveAvatarUrlForUser(anchorId, u.Avatar)
	}
	return item
}
