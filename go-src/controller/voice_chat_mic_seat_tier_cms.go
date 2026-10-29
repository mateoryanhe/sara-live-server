package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/voicechatmicseattierdto"
	"xr-game-server/module/voicechatmicseattier"
)

const VoiceChatMicSeatTierCMSUrl = "/voiceChatMicSeatTier"

type VoiceChatMicSeatTierCMSController struct{}

func initVoiceChatMicSeatTierCMSController() {
	httpserver.RegCMS(VoiceChatMicSeatTierCMSUrl, &VoiceChatMicSeatTierCMSController{})
}

func (c *VoiceChatMicSeatTierCMSController) VoiceChatMicSeatTierList(ctx context.Context, req *voicechatmicseattierdto.VoiceChatMicSeatTierListReq) (*httpserver.CMSQueryResp, error) {
	return voicechatmicseattier.ListVoiceChatMicSeatTiers(ctx, req)
}

func (c *VoiceChatMicSeatTierCMSController) CreateVoiceChatMicSeatTier(ctx context.Context, req *voicechatmicseattierdto.CreateVoiceChatMicSeatTierReq) (*voicechatmicseattierdto.CreateVoiceChatMicSeatTierRes, error) {
	return voicechatmicseattier.CreateVoiceChatMicSeatTier(ctx, req)
}

func (c *VoiceChatMicSeatTierCMSController) UpdateVoiceChatMicSeatTier(ctx context.Context, req *voicechatmicseattierdto.UpdateVoiceChatMicSeatTierReq) (*voicechatmicseattierdto.UpdateVoiceChatMicSeatTierRes, error) {
	return voicechatmicseattier.UpdateVoiceChatMicSeatTier(ctx, req)
}

func (c *VoiceChatMicSeatTierCMSController) DeleteVoiceChatMicSeatTier(ctx context.Context, req *voicechatmicseattierdto.DeleteVoiceChatMicSeatTierReq) (*voicechatmicseattierdto.DeleteVoiceChatMicSeatTierRes, error) {
	return voicechatmicseattier.DeleteVoiceChatMicSeatTier(ctx, req)
}

func (c *VoiceChatMicSeatTierCMSController) OnShelfVoiceChatMicSeatTier(ctx context.Context, req *voicechatmicseattierdto.OnShelfVoiceChatMicSeatTierReq) (*voicechatmicseattierdto.OnShelfVoiceChatMicSeatTierRes, error) {
	return voicechatmicseattier.OnShelfVoiceChatMicSeatTier(ctx, req)
}

func (c *VoiceChatMicSeatTierCMSController) OffShelfVoiceChatMicSeatTier(ctx context.Context, req *voicechatmicseattierdto.OffShelfVoiceChatMicSeatTierReq) (*voicechatmicseattierdto.OffShelfVoiceChatMicSeatTierRes, error) {
	return voicechatmicseattier.OffShelfVoiceChatMicSeatTier(ctx, req)
}
