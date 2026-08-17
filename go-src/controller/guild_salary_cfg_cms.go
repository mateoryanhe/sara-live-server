package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/guildsalarycfgdto"
	"xr-game-server/module/guildsalarycfg"
)

const GuildSalaryCfgCMSUrl = "/guildSalaryCfg"

type GuildSalaryCfgCMSController struct{}

func initGuildSalaryCfgCMSController() {
	httpserver.RegCMS(GuildSalaryCfgCMSUrl, &GuildSalaryCfgCMSController{})
}

func (c *GuildSalaryCfgCMSController) GuildSalaryCfgList(ctx context.Context, req *guildsalarycfgdto.GuildSalaryCfgListReq) (*httpserver.CMSQueryResp, error) {
	return guildsalarycfg.GetList(ctx, req)
}

func (c *GuildSalaryCfgCMSController) CreateGuildSalaryCfg(ctx context.Context, req *guildsalarycfgdto.CreateGuildSalaryCfgReq) (*guildsalarycfgdto.CreateGuildSalaryCfgRes, error) {
	return guildsalarycfg.Create(ctx, req)
}

func (c *GuildSalaryCfgCMSController) UpdateGuildSalaryCfg(ctx context.Context, req *guildsalarycfgdto.UpdateGuildSalaryCfgReq) (*guildsalarycfgdto.UpdateGuildSalaryCfgRes, error) {
	return guildsalarycfg.Update(ctx, req)
}

func (c *GuildSalaryCfgCMSController) DeleteGuildSalaryCfg(ctx context.Context, req *guildsalarycfgdto.DeleteGuildSalaryCfgReq) (*guildsalarycfgdto.DeleteGuildSalaryCfgRes, error) {
	return guildsalarycfg.Delete(ctx, req)
}
