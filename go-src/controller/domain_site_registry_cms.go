package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/core/httpserver"
	"xr-game-server/dto/domainsitedto"
	"xr-game-server/module/domainsite"
)

const DomainSiteRegistryCMSURL = "/domainSite"

type DomainSiteRegistryCMSController struct{}

func initDomainSiteRegistryCMSController() {
	httpserver.RegCMSHandler(DomainSiteRegistryCMSURL, "/saveExtraDomainMappingUpload", handleSaveExtraDomainMappingUpload)
	httpserver.RegCMSHandler(DomainSiteRegistryCMSURL, "/saveHostTlsMappingUpload", handleSaveHostTlsMappingUpload)
	httpserver.RegCMS(DomainSiteRegistryCMSURL, &DomainSiteRegistryCMSController{})
}

func handleSaveExtraDomainMappingUpload(r *ghttp.Request) {
	res, err := domainsite.SaveExtraDomainMappingFromRequest(r)
	if err != nil {
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func handleSaveHostTlsMappingUpload(r *ghttp.Request) {
	res, err := domainsite.SaveHostTlsMappingFromRequest(r)
	if err != nil {
		r.SetError(err)
		return
	}
	httpserver.SetHandlerResponseData(r, res)
}

func (c *DomainSiteRegistryCMSController) RefreshDomainSiteRegistry(ctx context.Context, req *domainsitedto.RefreshDomainSiteRegistryReq) (*domainsitedto.RefreshDomainSiteRegistryRes, error) {
	return domainsite.RefreshRegistry(ctx, req)
}

func (c *DomainSiteRegistryCMSController) GetCMSDomainSiteMapping(ctx context.Context, req *domainsitedto.GetCMSDomainSiteMappingReq) (*domainsitedto.GetCMSDomainSiteMappingRes, error) {
	return domainsite.GetCMSDomainSiteMapping(ctx, req)
}

func (c *DomainSiteRegistryCMSController) SaveCMSDomainSiteMapping(ctx context.Context, req *domainsitedto.SaveCMSDomainSiteMappingReq) (*domainsitedto.SaveCMSDomainSiteMappingRes, error) {
	return domainsite.SaveCMSDomainSiteMapping(ctx, req)
}

func (c *DomainSiteRegistryCMSController) ListExtraDomainMappings(ctx context.Context, req *domainsitedto.ListExtraDomainMappingsReq) (*domainsitedto.ListExtraDomainMappingsRes, error) {
	return domainsite.ListExtraDomainMappings(ctx, req)
}

func (c *DomainSiteRegistryCMSController) SaveExtraDomainMapping(ctx context.Context, req *domainsitedto.SaveExtraDomainMappingReq) (*domainsitedto.SaveExtraDomainMappingRes, error) {
	return domainsite.SaveExtraDomainMapping(ctx, req)
}

func (c *DomainSiteRegistryCMSController) DeleteExtraDomainMapping(ctx context.Context, req *domainsitedto.DeleteExtraDomainMappingReq) (*domainsitedto.DeleteExtraDomainMappingRes, error) {
	return domainsite.DeleteExtraDomainMapping(ctx, req)
}

func (c *DomainSiteRegistryCMSController) ListHostTlsMappings(ctx context.Context, req *domainsitedto.ListHostTlsMappingsReq) (*domainsitedto.ListHostTlsMappingsRes, error) {
	return domainsite.ListHostTlsMappings(ctx, req)
}

func (c *DomainSiteRegistryCMSController) DeleteHostTlsMapping(ctx context.Context, req *domainsitedto.DeleteHostTlsMappingReq) (*domainsitedto.DeleteHostTlsMappingRes, error) {
	return domainsite.DeleteHostTlsMapping(ctx, req)
}
