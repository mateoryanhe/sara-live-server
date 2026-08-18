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

// OffShelfGuildList 获取已下架工会列表(垃圾库)
func (c *GuildController) OffShelfGuildList(ctx context.Context, req *guilddto.OffShelfGuildListReq) (res *httpserver.CMSQueryResp, err error) {
	return guild.GetOffShelfGuildList(ctx, req)
}

// CreateGuild 创建直播工会
func (c *GuildController) CreateGuild(ctx context.Context, req *guilddto.CreateGuildReq) (res *guilddto.CreateGuildRes, err error) {
	return guild.CreateGuild(ctx, req)
}

// UpdateGuild 更新直播工会
func (c *GuildController) UpdateGuild(ctx context.Context, req *guilddto.UpdateGuildReq) (res *guilddto.UpdateGuildRes, err error) {
	return guild.UpdateGuild(ctx, req)
}

// DeleteGuild 下架直播工会
func (c *GuildController) DeleteGuild(ctx context.Context, req *guilddto.DeleteGuildReq) (res *guilddto.DeleteGuildRes, err error) {
	return guild.DeleteGuild(ctx, req)
}

// OnShelfGuild 上架直播工会
func (c *GuildController) OnShelfGuild(ctx context.Context, req *guilddto.OnShelfGuildReq) (res *guilddto.OnShelfGuildRes, err error) {
	return guild.OnShelfGuild(ctx, req)
}

// GetMyGuildProfile 获取当前 CMS 用户关联工会基础信息
func (c *GuildController) GetMyGuildProfile(ctx context.Context, req *guilddto.GetMyGuildProfileReq) (res *guilddto.GetMyGuildProfileRes, err error) {
	return guild.GetMyGuildProfile(ctx, req)
}

// UpdateMyGuildProfile 更新当前 CMS 用户关联工会基础信息
func (c *GuildController) UpdateMyGuildProfile(ctx context.Context, req *guilddto.UpdateMyGuildProfileReq) (res *guilddto.UpdateMyGuildProfileRes, err error) {
	return guild.UpdateMyGuildProfile(ctx, req)
}

// CMSMyGuildAnchorIncomeSettlementLogList 查询当前 CMS 用户名下工会的主播结算流水
func (c *GuildController) CMSMyGuildAnchorIncomeSettlementLogList(ctx context.Context, req *guilddto.CMSMyGuildAnchorIncomeSettlementLogListReq) (*httpserver.CMSQueryResp, error) {
	return guild.GetMyGuildAnchorIncomeSettlementLogList(ctx, req)
}

// ImportGuildAnchors CSV 导入工会主播
func (c *GuildController) ImportGuildAnchors(ctx context.Context, req *guilddto.ImportGuildAnchorsReq) (res *guilddto.ImportGuildAnchorsRes, err error) {
	return guild.ImportGuildAnchors(ctx, req)
}

// GetGuildDetail CMS获取工会详情收益
func (c *GuildController) GetGuildDetail(ctx context.Context, req *guilddto.GetGuildDetailReq) (res *guilddto.GetGuildDetailRes, err error) {
	return guild.QueryGuildDetail(ctx, req)
}

// GetGuildIncomeArchives CMS获取工会下架归档
func (c *GuildController) GetGuildIncomeArchives(ctx context.Context, req *guilddto.GetGuildIncomeArchivesReq) (res *guilddto.GetGuildIncomeArchivesRes, err error) {
	return guild.QueryGuildIncomeArchives(ctx, req)
}
