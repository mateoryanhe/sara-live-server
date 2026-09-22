package domainsite

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/core/cfg"
	"xr-game-server/dao/domainsitedao"
	"xr-game-server/dto/domainsitedto"
	"xr-game-server/entity/sys"
)

func GetCMSDomainSiteMapping(_ context.Context, _ *domainsitedto.GetCMSDomainSiteMappingReq) (*domainsitedto.GetCMSDomainSiteMappingRes, error) {
	return &domainsitedto.GetCMSDomainSiteMappingRes{
		Mapping: currentCMSDomainSiteMapping(),
	}, nil
}

func SaveCMSDomainSiteMapping(ctx context.Context, req *domainsitedto.SaveCMSDomainSiteMappingReq) (*domainsitedto.SaveCMSDomainSiteMappingRes, error) {
	row, err := SaveStaticSiteMapping(
		ctx,
		domainsitedto.CMSSiteKey,
		domainsitedto.CMSURLPrefix,
		strings.TrimSpace(req.Domain),
		strings.TrimSpace(req.Root),
	)
	if err != nil {
		return nil, err
	}
	return &domainsitedto.SaveCMSDomainSiteMappingRes{
		Success: true,
		Mapping: cmsDomainSiteMappingItem(row),
	}, nil
}

func currentCMSDomainSiteMapping() *domainsitedto.CMSDomainSiteMappingItem {
	if row := domainsitedao.GetMappingBySiteKey(domainsitedto.CMSSiteKey); row != nil {
		return cmsDomainSiteMappingItem(row)
	}
	root := strings.TrimSpace(cfg.GetStaticPathRoot(domainsitedto.CMSURLPrefix))
	if root == "" {
		root = domainsitedto.DefaultCMSRoot
	}
	return &domainsitedto.CMSDomainSiteMappingItem{
		ID:        "0",
		Domain:    cfg.GetStaticSiteDomain(domainsitedto.CMSURLPrefix),
		URLPrefix: domainsitedto.CMSURLPrefix,
		Root:      root,
	}
}

func cmsDomainSiteMappingItem(row *entity.DomainSiteMapping) *domainsitedto.CMSDomainSiteMappingItem {
	if row == nil {
		return nil
	}
	return &domainsitedto.CMSDomainSiteMappingItem{
		ID:        strconv.FormatUint(row.ID, 10),
		Domain:    row.Domain,
		URLPrefix: domainsitedto.CMSURLPrefix,
		Root:      row.Root,
		UpdatedAt: formatDomainSiteTime(row.UpdatedAt),
	}
}

func formatDomainSiteTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format("2006-01-02 15:04:05")
}
