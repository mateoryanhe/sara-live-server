package guild

import (
	"context"
	"strconv"
	"strings"

	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/guilddto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/liveroom"
)

// ImportGuildAnchors CMS 按工会 CSV 导入普通/高级主播
func ImportGuildAnchors(ctx context.Context, req *guilddto.ImportGuildAnchorsReq) (*guilddto.ImportGuildAnchorsRes, error) {
	if req == nil || req.GuildId == 0 || len(req.Rows) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if req.AnchorType != entity.UserTypeAnchor && req.AnchorType != entity.UserTypeSeniorAnchor {
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
		cancelCode := strings.TrimSpace(row.CancelCode)
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

		failReason, nickname := importOneGuildAnchor(guild, userID, cancelCode, req.AnchorType)
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

func importOneGuildAnchor(guild *entity.LiveGuild, userID uint64, cancelCode string, anchorType uint8) (failReason int, nickname string) {
	if accountdao.GetAccountById(userID) == nil {
		return guilddto.ImportAnchorFailUserNotFound, ""
	}

	user := userinfodao.GetUserInfoByUserId(userID)
	if user != nil {
		nickname = user.Nickname
	}

	ext := userinfodao.GetUserExtByUserId(userID)
	storedCode := ""
	if ext != nil {
		storedCode = strings.TrimSpace(ext.CancelCode)
	}
	if storedCode == "" || storedCode != cancelCode {
		return guilddto.ImportAnchorFailCancelCodeMismatch, nickname
	}
	if !userinfodao.IsCancelCodeValid(ext) {
		return guilddto.ImportAnchorFailCancelCodeExpired, nickname
	}

	if user.GuildId != 0 {
		return guilddto.ImportAnchorFailAlreadyInGuild, nickname
	}
	if !entity.UserTypeIsAnchor(user.UserType) && user.UserType != entity.UserTypeNormal {
		return guilddto.ImportAnchorFailCannotSetAnchor, nickname
	}

	user.SetGuildId(guild.ID)
	_ = liveroom.SetUserAsAnchorWithTimezone(userID, guild.ID, anchorType)
	return 0, nickname
}
