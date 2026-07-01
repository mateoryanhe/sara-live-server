package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/entryeffectdto"
	"xr-game-server/module/entryeffect"
)

const EntryEffectAppUrl = "/entryEffect"

type EntryEffectAppController struct{}

func initEntryEffectAppController() {
	httpserver.RegAPI(EntryEffectAppUrl, &EntryEffectAppController{})
}

func (c *EntryEffectAppController) AppEntryEffectList(ctx context.Context, req *entryeffectdto.AppEntryEffectListReq) (*entryeffectdto.AppEntryEffectListRes, error) {
	return entryeffect.GetAppEntryEffectList(ctx, req)
}
