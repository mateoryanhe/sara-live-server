package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/aliyuntextmoderationdto"
	"xr-game-server/module/aliyunmoderation"
)

const AliyunTextModerationCMSUrl = "/textModeration"

type AliyunTextModerationCMSController struct{}

func initAliyunTextModerationCMSController() {
	httpserver.RegCMS(AliyunTextModerationCMSUrl, &AliyunTextModerationCMSController{})
}

func (c *AliyunTextModerationCMSController) GetTextModerationCfg(ctx context.Context, req *aliyuntextmoderationdto.GetCfgReq) (*aliyuntextmoderationdto.GetCfgRes, error) {
	return aliyunmoderation.GetCfg(ctx, req)
}

func (c *AliyunTextModerationCMSController) SaveTextModerationCfg(ctx context.Context, req *aliyuntextmoderationdto.SaveCfgReq) (*aliyuntextmoderationdto.SaveCfgRes, error) {
	return aliyunmoderation.SaveCfg(ctx, req)
}
