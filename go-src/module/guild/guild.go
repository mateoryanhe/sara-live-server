package guild

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/guilddao"
	"xr-game-server/dto/guilddto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

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

	guild := entity.LiveGuild{
		Name:        req.Name,
		LeaderId:    req.LeaderId,
		Contact:     req.Contact,
		Description: req.Description,
		Status:      req.Status,
	}

	if err = guilddao.CreateGuild(&guild); err != nil {
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
	guild.Contact = req.Contact
	guild.Description = req.Description
	guild.Status = req.Status

	if err = guilddao.UpdateGuild(guild); err != nil {
		return nil, err
	}
	guilddao.RemoveGuildCache(req.ID)

	return &guilddto.UpdateGuildRes{Success: true}, nil
}

// DeleteGuild 删除直播工会
func DeleteGuild(ctx context.Context, req *guilddto.DeleteGuildReq) (res *guilddto.DeleteGuildRes, err error) {
	if err = guilddao.DeleteGuild(req.ID); err != nil {
		return nil, err
	}
	guilddao.RemoveGuildCache(req.ID)
	return &guilddto.DeleteGuildRes{Success: true}, nil
}

// GetGuild App端获取工会信息(走缓存)
func GetGuild(ctx context.Context, req *guilddto.GetGuildReq) (res *guilddto.GetGuildRes, err error) {
	guild := guilddao.GetGuildByIdCached(req.ID)
	if guild == nil {
		return nil, errercode.CreateCode(errercode.GuildNonExist)
	}
	return &guilddto.GetGuildRes{
		ID:          strconv.FormatUint(guild.ID, 10),
		Name:        guild.Name,
		LeaderId:    strconv.FormatUint(guild.LeaderId, 10),
		Contact:     guild.Contact,
		Description: guild.Description,
		Status:      guild.Status,
	}, nil
}
