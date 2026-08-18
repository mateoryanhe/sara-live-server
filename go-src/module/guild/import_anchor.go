package guild

import (
	"context"
	"strconv"

	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/guilddto"
	"xr-game-server/entity/live"
	userentity "xr-game-server/entity/user"
	"xr-game-server/errercode"
	"xr-game-server/module/liveroom"
)

// ImportGuildAnchors CMS 按工会 CSV 导入普通/高级主播(仅需用户ID)
func ImportGuildAnchors(ctx context.Context, req *guilddto.ImportGuildAnchorsReq) (*guilddto.ImportGuildAnchorsRes, error) {
	if req == nil || req.GuildId == 0 || len(req.Rows) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if req.AnchorType != userentity.UserTypeAnchor && req.AnchorType != userentity.UserTypeSeniorAnchor {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	guild := guilddao.GetGuildById(req.GuildId)
	if guild == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}

	res := &guilddto.ImportGuildAnchorsRes{
		Fails: make([]*guilddto.ImportGuildAnchorFailItem, 0),
	}
	seen := make(map[uint64]struct{}, len(req.Rows))
	successIDs := make([]uint64, 0)

	for _, row := range req.Rows {
		if row == nil {
			continue
		}
		userID := row.UserId
		if userID == 0 {
			res.Fails = append(res.Fails, &guilddto.ImportGuildAnchorFailItem{
				UserId: "0",
				Reason: guilddto.ImportAnchorFailUserNotFound,
			})
			continue
		}
		if _, ok := seen[userID]; ok {
			continue
		}
		seen[userID] = struct{}{}

		failReason, nickname := importOneGuildAnchor(guild, userID, req.AnchorType)
		if failReason > 0 {
			res.Fails = append(res.Fails, &guilddto.ImportGuildAnchorFailItem{
				UserId:   strconv.FormatUint(userID, 10),
				Nickname: nickname,
				Reason:   failReason,
			})
			continue
		}
		successIDs = append(successIDs, userID)
	}

	res.SuccessCount = len(successIDs)
	res.FailCount = len(res.Fails)
	if res.SuccessCount > 0 {
		liveroom.RefreshRoomListCache(ctx)
	}
	return res, nil
}

// JoinGuildAnchor CMS 单个用户加入工会(设为主播并绑定工会)
func JoinGuildAnchor(ctx context.Context, req *guilddto.SetAnchorGuildReq) (*guilddto.SetAnchorGuildRes, error) {
	if req == nil || req.GuildId == 0 || req.UserId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if req.AnchorType != userentity.UserTypeAnchor && req.AnchorType != userentity.UserTypeSeniorAnchor {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	guild := guilddao.GetGuildById(req.GuildId)
	if guild == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}
	failReason, nickname := importOneGuildAnchor(guild, req.UserId, req.AnchorType)
	if failReason > 0 {
		return &guilddto.SetAnchorGuildRes{
			Success:  false,
			Reason:   failReason,
			Nickname: nickname,
		}, nil
	}
	liveroom.RefreshRoomListCache(ctx)
	return &guilddto.SetAnchorGuildRes{Success: true}, nil
}

func importOneGuildAnchor(guild *entity.LiveGuild, userID uint64, anchorType uint8) (failReason int, nickname string) {
	if accountdao.GetAccountById(userID) == nil {
		return guilddto.ImportAnchorFailUserNotFound, ""
	}

	user := userinfodao.GetUserInfoByUserId(userID)
	if user == nil {
		return guilddto.ImportAnchorFailUserNotFound, ""
	}
	nickname = user.Nickname

	// 主播间缓存已有该ID,不可再导入当前工会
	if liveroomdao.GetRoomById(userID) != nil {
		return guilddto.ImportAnchorFailAlreadyHasLiveRoom, nickname
	}

	if liveroomdao.GetAnchorGuildId(userID) != 0 {
		return guilddto.ImportAnchorFailAlreadyInGuild, nickname
	}
	if !userentity.UserTypeIsAnchor(user.UserType) && user.UserType != userentity.UserTypeNormal {
		return guilddto.ImportAnchorFailCannotSetAnchor, nickname
	}

	_ = liveroom.SetUserAsAnchorIfNeeded(userID, anchorType)
	liveroom.EnsureAnchorRoom(userID, guild.ID)
	return 0, nickname
}
