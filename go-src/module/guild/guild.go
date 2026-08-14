package guild

import (
	"context"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/core/snowflake"
	"xr-game-server/dao/cmsuserdao"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dto/guilddto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/timezonecfg"
)

func genGuildId() uint64 {
	return snowflake.GetId()
}

// GetGuildList 获取直播工会列表
func GetGuildList(ctx context.Context, req *guilddto.GuildListReq) (res *httpserver.CMSQueryResp, err error) {
	total, guilds := guilddao.GetGuildList(req)
	return &httpserver.CMSQueryResp{
		Total: total,
		Data:  guilds,
	}, nil
}

// CreateGuild 创建直播工会
func CreateGuild(ctx context.Context, req *guilddto.CreateGuildReq) (res *guilddto.CreateGuildRes, err error) {
	if existing := guilddao.GetGuildByName(req.Name); existing != nil {
		return nil, errercode.CreateCode(errercode.GuildExist)
	}

	guild := entity.NewLiveGuild(
		genGuildId(),
		req.Name,
		req.LeaderId,
		resolveLeaderName(req.LeaderId),
		req.Description,
		req.Timezone,
	)
	if err = guilddao.CreateGuild(guild); err != nil {
		return nil, err
	}
	timezonecfg.EnsureCron(req.Timezone)

	return &guilddto.CreateGuildRes{ID: strconv.FormatUint(guild.ID, 10)}, nil
}

// UpdateGuild 更新直播工会
func UpdateGuild(ctx context.Context, req *guilddto.UpdateGuildReq) (res *guilddto.UpdateGuildRes, err error) {
	guild := guilddao.GetGuildById(req.ID)
	if guild == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}

	if existing := guilddao.GetGuildByName(req.Name); existing != nil && existing.ID != req.ID {
		return nil, errercode.CreateCode(errercode.GuildExist)
	}

	guild.Name = req.Name
	guild.LeaderId = req.LeaderId
	guild.LeaderName = resolveLeaderName(req.LeaderId)
	guild.Description = req.Description
	guild.Timezone = req.Timezone
	if err = guilddao.UpdateGuild(guild); err != nil {
		return nil, err
	}
	timezonecfg.EnsureCron(req.Timezone)

	return &guilddto.UpdateGuildRes{Success: true}, nil
}

// DeleteGuild 软删除直播工会
func DeleteGuild(ctx context.Context, req *guilddto.DeleteGuildReq) (res *guilddto.DeleteGuildRes, err error) {
	if guilddao.GetGuildById(req.ID) == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}
	if err = guilddao.DeleteGuild(req.ID); err != nil {
		return nil, err
	}
	return &guilddto.DeleteGuildRes{Success: true}, nil
}

func resolveLeaderName(leaderId uint64) string {
	if leaderId == 0 {
		return ""
	}
	user := cmsuserdao.GetCMSUserById(leaderId)
	if user == nil {
		return ""
	}
	return user.Name
}

// BatchUpdateGuildTimezone 批量更新工会时区
func BatchUpdateGuildTimezone(_ context.Context, req *guilddto.BatchUpdateGuildTimezoneReq) (res *guilddto.BatchUpdateGuildTimezoneRes, err error) {
	if err := guilddao.BatchUpdateGuildTimezone(req.GuildIds, req.Timezone); err != nil {
		return nil, err
	}
	timezonecfg.EnsureCron(req.Timezone)
	return &guilddto.BatchUpdateGuildTimezoneRes{Success: true}, nil
}
