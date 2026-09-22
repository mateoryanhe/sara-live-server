package domainsite

import (
	"context"
	"fmt"

	"xr-game-server/core/httpserver"
	"xr-game-server/dto/domainsitedto"
)

func RefreshRegistry(ctx context.Context, req *domainsitedto.RefreshDomainSiteRegistryReq) (*domainsitedto.RefreshDomainSiteRegistryRes, error) {
	var (
		active []httpserver.DomainSiteRegistration
		err    error
		source = "config"
	)
	if req != nil && req.Sites != nil {
		source = "request"
		registrations := make([]httpserver.DomainSiteRegistration, 0, len(req.Sites))
		for index, site := range req.Sites {
			if site == nil {
				return nil, fmt.Errorf("第 %d 项域名映射不能为空", index+1)
			}
			registrations = append(registrations, httpserver.DomainSiteRegistration{
				Domain:   site.Domain,
				Root:     site.Root,
				CertFile: site.CertFile,
				KeyFile:  site.KeyFile,
			})
		}
		active, err = httpserver.RefreshDomainSiteRegistry(ctx, registrations)
	} else {
		active, err = httpserver.RefreshDomainSiteRegistryFromConfig(ctx)
	}
	if err != nil {
		return nil, err
	}

	sites := make([]*domainsitedto.DomainSiteItem, 0, len(active))
	for _, site := range active {
		sites = append(sites, &domainsitedto.DomainSiteItem{
			Domain:   site.Domain,
			Root:     site.Root,
			CertFile: site.CertFile,
			KeyFile:  site.KeyFile,
		})
	}
	return &domainsitedto.RefreshDomainSiteRegistryRes{
		Count:  len(sites),
		Source: source,
		Sites:  sites,
	}, nil
}
