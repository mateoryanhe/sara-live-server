package guild

import (
	"context"
	"strconv"
	"time"

	"xr-game-server/dao/cmsuserdao"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dto/guilddto"
	"xr-game-server/errercode"
)

const (
	guildVisibilityTimeLayout = "2006-01-02 15:04:05"
	batchGuildVisibilityMax   = 500
)

func ensureGuildExists(guildId uint64) error {
	if guildId == 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if guilddao.GetGuildByIdFromDB(guildId) == nil {
		return errercode.CreateCode(errercode.GuildNonExist)
	}
	return nil
}

// GetGuildVisibilityList 查询工会已授权可见的 CMS 用户
func GetGuildVisibilityList(_ context.Context, req *guilddto.GuildVisibilityListReq) (*guilddto.GuildVisibilityListRes, error) {
	if err := ensureGuildExists(req.GuildId); err != nil {
		return nil, err
	}
	rows := guilddao.ListVisibilitiesByGuildId(req.GuildId)
	guildName := ""
	if g := guilddao.GetGuildByIdFromDB(req.GuildId); g != nil {
		guildName = g.Name
	}
	list := make([]*guilddto.GuildVisibilityItem, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		item := &guilddto.GuildVisibilityItem{
			Id:        strconv.FormatUint(row.ID, 10),
			GuildId:   strconv.FormatUint(row.GuildId, 10),
			GuildName: guildName,
			CmsUserId: strconv.FormatUint(row.CmsUserId, 10),
			CreatedAt: formatGuildVisibilityTime(row.CreatedAt),
		}
		if u := cmsuserdao.GetCMSUserById(row.CmsUserId); u != nil {
			item.CmsUserName = u.Name
		}
		list = append(list, item)
	}
	return &guilddto.GuildVisibilityListRes{List: list}, nil
}

// GetGuildVisibilityByUserList 查询指定 CMS 用户已授权的工会
func GetGuildVisibilityByUserList(_ context.Context, req *guilddto.GuildVisibilityByUserListReq) (*guilddto.GuildVisibilityByUserListRes, error) {
	if cmsuserdao.GetCMSUserById(req.CmsUserId) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	rows := guilddao.ListVisibilitiesByCmsUserId(req.CmsUserId)
	guildIds := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row != nil && row.GuildId > 0 {
			guildIds = append(guildIds, row.GuildId)
		}
	}
	nameMap := guilddao.GetNameMapByIds(guildIds)
	userName := ""
	if u := cmsuserdao.GetCMSUserById(req.CmsUserId); u != nil {
		userName = u.Name
	}
	list := make([]*guilddto.GuildVisibilityItem, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ID == 0 {
			continue
		}
		item := &guilddto.GuildVisibilityItem{
			Id:          strconv.FormatUint(row.ID, 10),
			GuildId:     strconv.FormatUint(row.GuildId, 10),
			CmsUserId:   strconv.FormatUint(row.CmsUserId, 10),
			CmsUserName: userName,
			CreatedAt:   formatGuildVisibilityTime(row.CreatedAt),
		}
		if nameMap != nil {
			item.GuildName = nameMap[row.GuildId]
		}
		list = append(list, item)
	}
	return &guilddto.GuildVisibilityByUserListRes{List: list}, nil
}

// GrantGuildVisibility 授权 CMS 用户可见某工会
func GrantGuildVisibility(_ context.Context, req *guilddto.GrantGuildVisibilityReq) (*guilddto.GrantGuildVisibilityRes, error) {
	if err := ensureGuildExists(req.GuildId); err != nil {
		return nil, err
	}
	if cmsuserdao.GetCMSUserById(req.CmsUserId) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := guilddao.AddGuildVisibility(req.GuildId, req.CmsUserId); err != nil {
		return nil, err
	}
	return &guilddto.GrantGuildVisibilityRes{Success: true}, nil
}

// BatchGrantGuildVisibility 批量授权 CMS 用户可见多个工会
func BatchGrantGuildVisibility(_ context.Context, req *guilddto.BatchGrantGuildVisibilityReq) (*guilddto.BatchGrantGuildVisibilityRes, error) {
	if cmsuserdao.GetCMSUserById(req.CmsUserId) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if len(req.GuildIds) == 0 || len(req.GuildIds) > batchGuildVisibilityMax {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	granted := 0
	seen := make(map[uint64]struct{}, len(req.GuildIds))
	for _, guildId := range req.GuildIds {
		if guildId == 0 {
			continue
		}
		if _, ok := seen[guildId]; ok {
			continue
		}
		seen[guildId] = struct{}{}
		if guilddao.GetGuildByIdFromDB(guildId) == nil {
			continue
		}
		if guilddao.HasGuildVisibility(guildId, req.CmsUserId) {
			continue
		}
		if err := guilddao.AddGuildVisibility(guildId, req.CmsUserId); err != nil {
			return nil, err
		}
		granted++
	}
	return &guilddto.BatchGrantGuildVisibilityRes{Success: true, GrantedCount: granted}, nil
}

// RevokeGuildVisibility 撤销 CMS 用户对某工会的可见性
func RevokeGuildVisibility(_ context.Context, req *guilddto.RevokeGuildVisibilityReq) (*guilddto.RevokeGuildVisibilityRes, error) {
	if err := ensureGuildExists(req.GuildId); err != nil {
		return nil, err
	}
	if !guilddao.HasGuildVisibility(req.GuildId, req.CmsUserId) {
		return &guilddto.RevokeGuildVisibilityRes{Success: true}, nil
	}
	if err := guilddao.RemoveGuildVisibility(req.GuildId, req.CmsUserId); err != nil {
		return nil, err
	}
	return &guilddto.RevokeGuildVisibilityRes{Success: true}, nil
}

// BatchRevokeGuildVisibility 批量撤销 CMS 用户对多个工会的可见性
func BatchRevokeGuildVisibility(_ context.Context, req *guilddto.BatchRevokeGuildVisibilityReq) (*guilddto.BatchRevokeGuildVisibilityRes, error) {
	if len(req.GuildIds) == 0 || len(req.GuildIds) > batchGuildVisibilityMax {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	revoked := 0
	seen := make(map[uint64]struct{}, len(req.GuildIds))
	for _, guildId := range req.GuildIds {
		if guildId == 0 {
			continue
		}
		if _, ok := seen[guildId]; ok {
			continue
		}
		seen[guildId] = struct{}{}
		if !guilddao.HasGuildVisibility(guildId, req.CmsUserId) {
			continue
		}
		if err := guilddao.RemoveGuildVisibility(guildId, req.CmsUserId); err != nil {
			return nil, err
		}
		revoked++
	}
	return &guilddto.BatchRevokeGuildVisibilityRes{Success: true, RevokedCount: revoked}, nil
}

func formatGuildVisibilityTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(guildVisibilityTimeLayout)
}
