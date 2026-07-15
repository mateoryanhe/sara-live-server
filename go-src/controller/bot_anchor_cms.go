package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/botanchordto"
	"xr-game-server/module/botanchor"
)

const (
	BotAnchorUrl = "/botAnchor"
)

type BotAnchorController struct {
}

func initBotAnchorController() {
	httpserver.RegCMS(BotAnchorUrl, &BotAnchorController{})
}

func (c *BotAnchorController) GetBotAnchorList(ctx context.Context, req *botanchordto.QueryBotAnchorListReq) (res *httpserver.CMSQueryResp, err error) {
	return botanchor.QueryBotAnchorList(ctx, req)
}

func (c *BotAnchorController) CreateBotAnchor(ctx context.Context, req *botanchordto.CreateBotAnchorReq) (res *botanchordto.CreateBotAnchorRes, err error) {
	return botanchor.CreateBotAnchor(ctx, req)
}

func (c *BotAnchorController) UpdateBotAnchor(ctx context.Context, req *botanchordto.UpdateBotAnchorReq) (res *botanchordto.UpdateBotAnchorRes, err error) {
	return botanchor.UpdateBotAnchor(ctx, req)
}

func (c *BotAnchorController) SetBotAnchorStatus(ctx context.Context, req *botanchordto.SetBotAnchorStatusReq) (res *botanchordto.SetBotAnchorStatusRes, err error) {
	return botanchor.SetBotAnchorStatus(ctx, req)
}
