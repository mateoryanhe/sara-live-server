package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/voicechatmicseattierdto"
	"xr-game-server/module/voicechatmicseattier"
)

const VoiceChatMicSeatTierAppUrl = "/voiceChatMicSeatTier"

type VoiceChatMicSeatTierAppController struct{}

func initVoiceChatMicSeatTierAppController() {
	httpserver.RegAPI(VoiceChatMicSeatTierAppUrl, &VoiceChatMicSeatTierAppController{})
}

func (c *VoiceChatMicSeatTierAppController) AppTierList(ctx context.Context, req *voicechatmicseattierdto.AppTierListReq) (*voicechatmicseattierdto.AppTierListRes, error) {
	return voicechatmicseattier.AppTierList(ctx, req)
}
