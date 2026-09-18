package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/staticcachecfgdto"
	"xr-game-server/module/staticcachecfg"
)

const StaticCacheCfgCMSURL = "/staticCacheCfg"

type StaticCacheCfgCMSController struct{}

func initStaticCacheCfgCMSController() {
	httpserver.RegCMS(StaticCacheCfgCMSURL, &StaticCacheCfgCMSController{})
}

func (c *StaticCacheCfgCMSController) StaticCacheRuleList(ctx context.Context, req *staticcachecfgdto.StaticCacheRuleListReq) (*httpserver.CMSQueryResp, error) {
	return staticcachecfg.GetList(ctx, req)
}

func (c *StaticCacheCfgCMSController) CreateStaticCacheRule(ctx context.Context, req *staticcachecfgdto.CreateStaticCacheRuleReq) (*staticcachecfgdto.CreateStaticCacheRuleRes, error) {
	return staticcachecfg.Create(ctx, req)
}

func (c *StaticCacheCfgCMSController) UpdateStaticCacheRule(ctx context.Context, req *staticcachecfgdto.UpdateStaticCacheRuleReq) (*staticcachecfgdto.UpdateStaticCacheRuleRes, error) {
	return staticcachecfg.Update(ctx, req)
}

func (c *StaticCacheCfgCMSController) DeleteStaticCacheRule(ctx context.Context, req *staticcachecfgdto.DeleteStaticCacheRuleReq) (*staticcachecfgdto.DeleteStaticCacheRuleRes, error) {
	return staticcachecfg.Delete(ctx, req)
}
