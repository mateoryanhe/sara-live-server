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

// GetMyGuildAnchorList 获取当前 CMS 用户管理的工会名下主播列表
func (c *GuildController) GetMyGuildAnchorList(ctx context.Context, req *guilddto.GetMyGuildAnchorListReq) (*httpserver.CMSQueryResp, error) {
	return guild.GetMyGuildAnchorList(ctx, req)
}

// GetMyGuildAnchorDailyEffectiveLiveList 获取当前 CMS 用户管理的指定工会名下主播每日流水
func (c *GuildController) GetMyGuildAnchorDailyEffectiveLiveList(ctx context.Context, req *guilddto.GetMyGuildAnchorDailyEffectiveLiveListReq) (*httpserver.CMSQueryResp, error) {
	return guild.GetMyGuildAnchorDailyEffectiveLiveList(ctx, req)
}

// GetMyOwnedGuildAnchorDailyEffectiveLiveList 获取当前 CMS 用户管理的工会名下主播每日流水
func (c *GuildController) GetMyOwnedGuildAnchorDailyEffectiveLiveList(ctx context.Context, req *guilddto.CMSMyGuildAnchorDailyEffectiveLiveListReq) (*httpserver.CMSQueryResp, error) {
	return guild.GetMyOwnedGuildAnchorDailyEffectiveLiveList(ctx, req)
}

// GetMyOwnedGuildAnchorList 获取当前 CMS 用户管理的全部工会名下主播
func (c *GuildController) GetMyOwnedGuildAnchorList(ctx context.Context, req *guilddto.GetMyOwnedGuildAnchorListReq) (*httpserver.CMSQueryResp, error) {
	return guild.GetMyOwnedGuildAnchorList(ctx, req)
}

// CMSGuildAnchorIncomeSettlementLogList CMS查询指定工会名下主播结算流水
func (c *GuildController) CMSGuildAnchorIncomeSettlementLogList(ctx context.Context, req *guilddto.CMSGuildAnchorIncomeSettlementLogListReq) (*httpserver.CMSQueryResp, error) {
	return guild.GetGuildAnchorIncomeSettlementLogList(ctx, req)
}

// CMSMyGuildAnchorIncomeSettlementLogList 查询当前 CMS 用户名下工会的主播结算流水
func (c *GuildController) CMSMyGuildAnchorIncomeSettlementLogList(ctx context.Context, req *guilddto.CMSMyGuildAnchorIncomeSettlementLogListReq) (*httpserver.CMSQueryResp, error) {
	return guild.GetMyGuildAnchorIncomeSettlementLogList(ctx, req)
}

// ImportGuildAnchors CSV 导入工会主播
func (c *GuildController) ImportGuildAnchors(ctx context.Context, req *guilddto.ImportGuildAnchorsReq) (res *guilddto.ImportGuildAnchorsRes, err error) {
	return guild.ImportGuildAnchors(ctx, req)
}

// JoinGuildAnchor CMS加入工会主播
func (c *GuildController) JoinGuildAnchor(ctx context.Context, req *guilddto.SetAnchorGuildReq) (res *guilddto.SetAnchorGuildRes, err error) {
	return guild.JoinGuildAnchor(ctx, req)
}

// SetGuildAnchorType CMS设置工会主播类型
func (c *GuildController) SetGuildAnchorType(ctx context.Context, req *guilddto.SetGuildAnchorTypeReq) (res *guilddto.SetGuildAnchorTypeRes, err error) {
	return guild.SetGuildAnchorType(ctx, req)
}

// GetGuildDetail CMS获取工会详情收益
func (c *GuildController) GetGuildDetail(ctx context.Context, req *guilddto.GetGuildDetailReq) (res *guilddto.GetGuildDetailRes, err error) {
	return guild.QueryGuildDetail(ctx, req)
}

// GetGuildIncomeArchives CMS获取工会下架归档
func (c *GuildController) GetGuildIncomeArchives(ctx context.Context, req *guilddto.GetGuildIncomeArchivesReq) (res *guilddto.GetGuildIncomeArchivesRes, err error) {
	return guild.QueryGuildIncomeArchives(ctx, req)
}

// GetGuildDailyEffectiveLiveList CMS查询工会每日流水
func (c *GuildController) GetGuildDailyEffectiveLiveList(ctx context.Context, req *guilddto.GetGuildDailyEffectiveLiveListReq) (*httpserver.CMSQueryResp, error) {
	return guild.QueryGuildDailyEffectiveLiveList(ctx, req)
}

// CMSGuildAnchorDailyEffectiveLiveList CMS查询工会名下主播每日流水
func (c *GuildController) CMSGuildAnchorDailyEffectiveLiveList(ctx context.Context, req *guilddto.CMSGuildAnchorDailyEffectiveLiveListReq) (*httpserver.CMSQueryResp, error) {
	return guild.GetGuildAnchorDailyEffectiveLiveList(ctx, req)
}

// GetGuildTransferInfo CMS获取工会收款/转账信息
func (c *GuildController) GetGuildTransferInfo(ctx context.Context, req *guilddto.GetGuildTransferInfoReq) (*guilddto.GetGuildTransferInfoRes, error) {
	return guild.GetGuildTransferInfo(ctx, req)
}

// SaveGuildTransferInfo CMS保存工会收款/转账信息(直写DB)
func (c *GuildController) SaveGuildTransferInfo(ctx context.Context, req *guilddto.SaveGuildTransferInfoReq) (*guilddto.SaveGuildTransferInfoRes, error) {
	return guild.SaveGuildTransferInfo(ctx, req)
}
