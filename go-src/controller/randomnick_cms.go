package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/randomnickdto"
	"xr-game-server/module/randomnick"
)

const RandomNicknameUrl = "/randomNickname"

type RandomNicknameController struct{}

func initRandomNicknameController() {
	httpserver.RegCMS(RandomNicknameUrl, &RandomNicknameController{})
}

func (c *RandomNicknameController) GetRandomNicknameCfg(ctx context.Context, req *randomnickdto.GetRandomNicknameCfgReq) (*randomnickdto.GetRandomNicknameCfgRes, error) {
	_ = req
	return randomnick.GetCMSCfgDTO(ctx)
}

func (c *RandomNicknameController) ImportRandomNicknames(ctx context.Context, req *randomnickdto.ImportRandomNicknamesReq) (*randomnickdto.ImportRandomNicknamesRes, error) {
	return randomnick.ImportNicknamesDTO(ctx, req)
}

func (c *RandomNicknameController) ClearRandomNicknames(ctx context.Context, req *randomnickdto.ClearRandomNicknamesReq) (*randomnickdto.ClearRandomNicknamesRes, error) {
	return randomnick.ClearNicknamesDTO(ctx, req)
}
