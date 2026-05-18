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

// ApplyJoin 申请加入工会
func (c *GuildAppController) ApplyJoin(ctx context.Context, req *guilddto.ApplyJoinGuildReq) (res *guilddto.ApplyJoinGuildRes, err error) {
	return guild.ApplyJoin(ctx, req)
}

// RejectMember 拒绝工会申请
func (c *GuildAppController) RejectMember(ctx context.Context, req *guilddto.RejectGuildMemberReq) (res *guilddto.RejectGuildMemberRes, err error) {
	return guild.RejectMember(ctx, req)
}

// KickMember 剔除工会成员
func (c *GuildAppController) KickMember(ctx context.Context, req *guilddto.KickGuildMemberReq) (res *guilddto.KickGuildMemberRes, err error) {
	return guild.KickMember(ctx, req)
}

// MemberList 获取工会成员列表(仅已加入)
func (c *GuildAppController) MemberList(ctx context.Context, req *guilddto.GuildMemberListReq) (res *guilddto.GuildMemberListRes, err error) {
	return guild.GetMemberList(ctx, req)
}

// PendingMemberList 获取工会待审批的入会申请列表(仅会长可调用)
func (c *GuildAppController) PendingMemberList(ctx context.Context, req *guilddto.PendingMemberListReq) (res *guilddto.PendingMemberListRes, err error) {
	return guild.GetPendingMemberList(ctx, req)
}
