package domainsitedao

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/sys"
)

func ListMappings() []*entity.DomainSiteMapping {
	rows := make([]*entity.DomainSiteMapping, 0)
	_ = g.DB().Model(string(entity.TbDomainSiteMapping)).Order("id asc").Scan(&rows)
	return rows
}

func GetMappingBySiteKey(siteKey string) *entity.DomainSiteMapping {
	var row entity.DomainSiteMapping
	if err := g.DB().Model(string(entity.TbDomainSiteMapping)).
		Where("site_key = ?", strings.TrimSpace(siteKey)).
		Limit(1).
		Scan(&row); err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

func SaveMapping(row *entity.DomainSiteMapping) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbDomainSiteMapping)).Save(row)
	return err
}

func DeleteMapping(id uint64) error {
	if id == 0 {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbDomainSiteMapping)).WherePri(id).Delete()
	return err
}
