package controller

import (
	"context"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/entryeffectdto"
	"xr-game-server/module/entryeffect"
)

const EntryEffectUrl = "/entryEffect"

type EntryEffectController struct{}

func initEntryEffectController() {
	httpserver.RegCMS(EntryEffectUrl, &EntryEffectController{})
}

func (c *EntryEffectController) EntryEffectList(ctx context.Context, req *entryeffectdto.EntryEffectListReq) (*httpserver.CMSQueryResp, error) {
	return entryeffect.GetEntryEffectList(ctx, req)
}

func (c *EntryEffectController) CreateEntryEffect(ctx context.Context, req *entryeffectdto.CreateEntryEffectReq) (*entryeffectdto.CreateEntryEffectRes, error) {
	return entryeffect.CreateEntryEffect(ctx, req)
}

func (c *EntryEffectController) UpdateEntryEffect(ctx context.Context, req *entryeffectdto.UpdateEntryEffectReq) (*entryeffectdto.UpdateEntryEffectRes, error) {
	return entryeffect.UpdateEntryEffect(ctx, req)
}

func (c *EntryEffectController) DeleteEntryEffect(ctx context.Context, req *entryeffectdto.DeleteEntryEffectReq) (*entryeffectdto.DeleteEntryEffectRes, error) {
	return entryeffect.DeleteEntryEffect(ctx, req)
}

func (c *EntryEffectController) OnShelfEntryEffect(ctx context.Context, req *entryeffectdto.OnShelfEntryEffectReq) (*entryeffectdto.OnShelfEntryEffectRes, error) {
	return entryeffect.OnShelfEntryEffect(ctx, req)
}

func (c *EntryEffectController) OffShelfEntryEffect(ctx context.Context, req *entryeffectdto.OffShelfEntryEffectReq) (*entryeffectdto.OffShelfEntryEffectRes, error) {
	return entryeffect.OffShelfEntryEffect(ctx, req)
}
