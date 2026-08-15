package guild

import (
	"context"
	"strconv"

	"xr-game-server/core/httpserver"
	"xr-game-server/core/snowflake"
	"xr-game-server/dao/cmsuserdao"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dto/guilddto"
	liveentity "xr-game-server/entity/live"
	"xr-game-server/errercode"
	"xr-game-server/module/liveroom"
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

// GetOffShelfGuildList 获取已下架工会列表(垃圾库,直查DB)
func GetOffShelfGuildList(ctx context.Context, req *guilddto.OffShelfGuildListReq) (res *httpserver.CMSQueryResp, err error) {
	total, guilds := guilddao.GetOffShelfGuildList(req)
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

	guild := liveentity.NewLiveGuild(
		genGuildId(),
		req.Name,
		req.LeaderId,
		resolveLeaderName(req.LeaderId),
		req.Description,
	)
	if err = guilddao.CreateGuild(guild); err != nil {
		return nil, err
	}

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
	if err = guilddao.UpdateGuild(guild); err != nil {
		return nil, err
	}

	return &guilddto.UpdateGuildRes{Success: true}, nil
}

// DeleteGuild 下架直播工会(同步下架工会下全部主播间)
func DeleteGuild(ctx context.Context, req *guilddto.DeleteGuildReq) (res *guilddto.DeleteGuildRes, err error) {
	if guilddao.GetGuildById(req.ID) == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}
	liveroom.OffShelfGuildLiveRooms(ctx, req.ID)
	if err = guilddao.DeleteGuild(req.ID); err != nil {
		return nil, err
	}
	return &guilddto.DeleteGuildRes{Success: true}, nil
}

// OnShelfGuild 上架直播工会(同步上架工会下全部主播间)
func OnShelfGuild(ctx context.Context, req *guilddto.OnShelfGuildReq) (res *guilddto.OnShelfGuildRes, err error) {
	guild := guilddao.GetGuildByIdFromDB(req.ID)
	if guild == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}
	if guild.Status == liveentity.LiveGuildStatusOnShelf {
		return &guilddto.OnShelfGuildRes{Success: true}, nil
	}
	// 上架前校验名称是否与已上架工会冲突
	if existing := guilddao.GetGuildByName(guild.Name); existing != nil && existing.ID != guild.ID {
		return nil, errercode.CreateCode(errercode.GuildExist)
	}
	if err = guilddao.OnShelfGuild(req.ID); err != nil {
		return nil, err
	}
	liveroom.OnShelfGuildLiveRooms(ctx, req.ID)
	return &guilddto.OnShelfGuildRes{Success: true}, nil
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
