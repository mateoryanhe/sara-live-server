package liveroom

import (
	"context"

	"xr-game-server/dao/guilddao"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/accountdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/timezonecfg"
)

// SetAnchor CMS 将用户设为主播(仅允许从未主播变为主播,不可回退)
func SetAnchor(ctx context.Context, req *accountdto.SetAnchorReq) (*accountdto.SetAnchorRes, error) {
	res, err := setUserAsAnchor(req.AccountId, entity.UserTypeAnchor)
	if err != nil {
		return nil, err
	}
	RefreshRoomListCache(ctx)
	return res, nil
}

// SetSeniorAnchor CMS 将用户设为高级主播(仅允许普通用户,不可回退)
func SetSeniorAnchor(ctx context.Context, req *accountdto.SetSeniorAnchorReq) (*accountdto.SetSeniorAnchorRes, error) {
	_, err := setUserAsAnchor(req.AccountId, entity.UserTypeSeniorAnchor)
	if err != nil {
		return nil, err
	}
	RefreshRoomListCache(ctx)
	return &accountdto.SetSeniorAnchorRes{Success: true}, nil
}

// BatchSetAnchor CMS 批量设普通主播
func BatchSetAnchor(ctx context.Context, req *accountdto.BatchSetAnchorReq) (*accountdto.BatchSetAnchorRes, error) {
	return batchSetAnchor(ctx, req.IDs, entity.UserTypeAnchor)
}

// BatchSetSeniorAnchor CMS 批量设高级主播
func BatchSetSeniorAnchor(ctx context.Context, req *accountdto.BatchSetSeniorAnchorReq) (*accountdto.BatchSetAnchorRes, error) {
	return batchSetAnchor(ctx, req.IDs, entity.UserTypeSeniorAnchor)
}

func batchSetAnchor(ctx context.Context, ids []uint64, userType uint8) (*accountdto.BatchSetAnchorRes, error) {
	res := &accountdto.BatchSetAnchorRes{}
	seen := make(map[uint64]struct{}, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		if err := setUserAsAnchorBatch(id, userType); err != nil {
			res.FailCount++
			res.FailIds = append(res.FailIds, id)
			continue
		}
		res.SuccessCount++
	}
	if res.SuccessCount > 0 {
		RefreshRoomListCache(ctx)
	}
	return res, nil
}

func setUserAsAnchor(accountId uint64, userType uint8) (*accountdto.SetAnchorRes, error) {
	if err := setUserAsAnchorBatch(accountId, userType); err != nil {
		return nil, err
	}
	return &accountdto.SetAnchorRes{Success: true}, nil
}

func setUserAsAnchorBatch(accountId uint64, userType uint8) error {
	user := userinfodao.GetUserInfoByUserId(accountId)
	if entity.UserTypeIsAnchor(user.UserType) {
		return nil
	}
	if user.UserType != entity.UserTypeNormal {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	user.SetUserType(userType)
	EnsureAnchorRoom(accountId, user.GuildId)
	return nil
}

// SetUserAsAnchorIfNeeded 将普通用户设为主播类型;已是主播则跳过
func SetUserAsAnchorIfNeeded(accountId uint64, userType uint8) error {
	return setUserAsAnchorBatch(accountId, userType)
}

// SetUserAsAnchorWithTimezone 将普通用户设为主播类型;已是主播则只更新时区;时区从工会获取
func SetUserAsAnchorWithTimezone(accountId, guildId uint64, userType uint8) error {
	user := userinfodao.GetUserInfoByUserId(accountId)
	if user == nil {
		return nil
	}
	if !entity.UserTypeIsAnchor(user.UserType) && user.UserType != entity.UserTypeNormal {
		return nil
	}
	// 从工会获取时区
	guild := guilddao.GetGuildById(guildId)
	timezone := 0
	if guild != nil {
		timezone = guild.Timezone
	}
	EnsureAnchorRoomWithTimezone(accountId, guildId, timezone)
	if entity.UserTypeIsAnchor(user.UserType) {
		return nil
	}
	user.SetUserType(userType)
	return nil
}

// BatchSetAnchorTimezone CMS批量设置主播时区(仅限工会ID=0的主播)
func BatchSetAnchorTimezone(ctx context.Context, req *accountdto.BatchSetAnchorTimezoneReq) (*accountdto.BatchSetAnchorTimezoneRes, error) {
	res := &accountdto.BatchSetAnchorTimezoneRes{}
	for _, anchorId := range req.AnchorIds {
		if err := setAnchorTimezone(anchorId, req.Timezone); err != nil {
			res.FailCount++
			res.FailIds = append(res.FailIds, anchorId)
			continue
		}
		res.SuccessCount++
	}
	if res.SuccessCount > 0 {
		timezonecfg.EnsureCron(req.Timezone)
		RefreshRoomListCache(ctx)
	}
	return res, nil
}

// setAnchorTimezone 设置单个主播时区(仅工会ID=0才可设置)
func setAnchorTimezone(anchorId uint64, timezone int) error {
	user := userinfodao.GetUserInfoByUserId(anchorId)
	if user == nil {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if !entity.UserTypeIsAnchor(user.UserType) {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if user.GuildId != 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	room := liveroomdao.GetRoomByAnchor(anchorId)
	if room == nil {
		return nil
	}
	room.SetTimezone(timezone)
	return nil
}

// ExitGuild CMS主播退出工会(将工会ID置为0)
func ExitGuild(ctx context.Context, req *accountdto.ExitGuildReq) (*accountdto.ExitGuildRes, error) {
	user := userinfodao.GetUserInfoByUserId(req.AnchorId)
	if user == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if !entity.UserTypeIsAnchor(user.UserType) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	user.SetGuildId(0)
	RefreshRoomListCache(ctx)
	return &accountdto.ExitGuildRes{Success: true}, nil
}
