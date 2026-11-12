package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/anchorvideopricetierdto"
	"xr-game-server/module/anchorvideopricetier"
)

const AnchorVideoPriceTierCMSUrl = "/anchorVideoPriceTier"

type AnchorVideoPriceTierCMSController struct{}

func initAnchorVideoPriceTierCMSController() {
	httpserver.RegCMS(AnchorVideoPriceTierCMSUrl, &AnchorVideoPriceTierCMSController{})
}

func (c *AnchorVideoPriceTierCMSController) AnchorVideoPriceTierList(ctx context.Context, req *anchorvideopricetierdto.PriceTierListReq) (*httpserver.CMSQueryResp, error) {
	return anchorvideopricetier.GetPriceTierList(ctx, req)
}

func (c *AnchorVideoPriceTierCMSController) CreateAnchorVideoPriceTier(ctx context.Context, req *anchorvideopricetierdto.CreatePriceTierReq) (*anchorvideopricetierdto.CreatePriceTierRes, error) {
	return anchorvideopricetier.CreatePriceTier(ctx, req)
}

func (c *AnchorVideoPriceTierCMSController) UpdateAnchorVideoPriceTier(ctx context.Context, req *anchorvideopricetierdto.UpdatePriceTierReq) (*anchorvideopricetierdto.UpdatePriceTierRes, error) {
	return anchorvideopricetier.UpdatePriceTier(ctx, req)
}

func (c *AnchorVideoPriceTierCMSController) DeleteAnchorVideoPriceTier(ctx context.Context, req *anchorvideopricetierdto.DeletePriceTierReq) (*anchorvideopricetierdto.DeletePriceTierRes, error) {
	return anchorvideopricetier.DeletePriceTier(ctx, req)
}

func (c *AnchorVideoPriceTierCMSController) OnShelfAnchorVideoPriceTier(ctx context.Context, req *anchorvideopricetierdto.OnShelfPriceTierReq) (*anchorvideopricetierdto.OnShelfPriceTierRes, error) {
	return anchorvideopricetier.OnShelfPriceTier(ctx, req)
}

func (c *AnchorVideoPriceTierCMSController) OffShelfAnchorVideoPriceTier(ctx context.Context, req *anchorvideopricetierdto.OffShelfPriceTierReq) (*anchorvideopricetierdto.OffShelfPriceTierRes, error) {
	return anchorvideopricetier.OffShelfPriceTier(ctx, req)
}
