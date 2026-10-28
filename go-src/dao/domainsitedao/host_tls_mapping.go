package domainsitedao

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/sys"
)

func ListHostTlsMappings() []*entity.HostTlsMapping {
	rows := make([]*entity.HostTlsMapping, 0)
	_ = g.DB().Model(string(entity.TbHostTlsMapping)).Order("id asc").Scan(&rows)
	return rows
}

func GetHostTlsMappingByDomain(domain string) *entity.HostTlsMapping {
	var row entity.HostTlsMapping
	if err := g.DB().Model(string(entity.TbHostTlsMapping)).
		Where("domain = ?", strings.TrimSpace(domain)).
		Limit(1).
		Scan(&row); err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

func GetHostTlsMappingByID(id uint64) *entity.HostTlsMapping {
	if id == 0 {
		return nil
	}
	var row entity.HostTlsMapping
	if err := g.DB().Model(string(entity.TbHostTlsMapping)).WherePri(id).Scan(&row); err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

func SaveHostTlsMapping(row *entity.HostTlsMapping) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbHostTlsMapping)).Save(row)
	return err
}

func DeleteHostTlsMapping(id uint64) error {
	if id == 0 {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbHostTlsMapping)).WherePri(id).Delete()
	return err
}
