package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/guilddto"
	"xr-game-server/module/guild"
)

const (
	GuildUrl = "/guild"
)

type GuildController struct {
}

func initGuildController() {
	httpserver.RegCMS(GuildUrl, &GuildController{})
}

// GuildList 获取直播工会列表
func (c *GuildController) GuildList(ctx context.Context, req *guilddto.GuildListReq) (res *httpserver.CMSQueryResp, err error) {
	return guild.GetGuildList(ctx, req)
}

// CreateGuild 创建直播工会
func (c *GuildController) CreateGuild(ctx context.Context, req *guilddto.CreateGuildReq) (res *guilddto.CreateGuildRes, err error) {
	return guild.CreateGuild(ctx, req)
}

// UpdateGuild 更新直播工会
func (c *GuildController) UpdateGuild(ctx context.Context, req *guilddto.UpdateGuildReq) (res *guilddto.UpdateGuildRes, err error) {
	return guild.UpdateGuild(ctx, req)
}

// DeleteGuild 删除直播工会
func (c *GuildController) DeleteGuild(ctx context.Context, req *guilddto.DeleteGuildReq) (res *guilddto.DeleteGuildRes, err error) {
	return guild.DeleteGuild(ctx, req)
}
