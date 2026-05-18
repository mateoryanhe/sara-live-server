package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/guilddto"
	"xr-game-server/module/guild"
)

const (
	GuildAppUrl = "/guild"
)

type GuildAppController struct {
}

func initGuildAppController() {
	httpserver.RegAPI(GuildAppUrl, &GuildAppController{})
}

// GetGuild App端获取工会信息
func (c *GuildAppController) GetGuild(ctx context.Context, req *guilddto.GetGuildReq) (res *guilddto.GetGuildRes, err error) {
	return guild.GetGuild(ctx, req)
}
