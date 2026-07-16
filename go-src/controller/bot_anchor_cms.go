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

func (c *BotAnchorController) StartBotAnchorLive(ctx context.Context, req *botanchordto.StartBotAnchorLiveReq) (res *botanchordto.StartBotAnchorLiveRes, err error) {
	return botanchor.StartBotAnchorLive(ctx, req)
}

func (c *BotAnchorController) StopBotAnchorLive(ctx context.Context, req *botanchordto.StopBotAnchorLiveReq) (res *botanchordto.StopBotAnchorLiveRes, err error) {
	return botanchor.StopBotAnchorLive(ctx, req)
}

func (c *BotAnchorController) BatchStartBotAnchorLive(ctx context.Context, req *botanchordto.BatchStartBotAnchorLiveReq) (res *botanchordto.BatchBotAnchorLiveRes, err error) {
	return botanchor.BatchStartBotAnchorLive(ctx, req)
}

func (c *BotAnchorController) BatchStopBotAnchorLive(ctx context.Context, req *botanchordto.BatchStopBotAnchorLiveReq) (res *botanchordto.BatchBotAnchorLiveRes, err error) {
	return botanchor.BatchStopBotAnchorLive(ctx, req)
}
