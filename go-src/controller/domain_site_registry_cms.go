package controller

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/domainsitedto"
	"xr-game-server/module/domainsite"
)

const DomainSiteRegistryCMSURL = "/domainSite"

type DomainSiteRegistryCMSController struct{}

func initDomainSiteRegistryCMSController() {
	httpserver.RegCMS(DomainSiteRegistryCMSURL, &DomainSiteRegistryCMSController{})
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
