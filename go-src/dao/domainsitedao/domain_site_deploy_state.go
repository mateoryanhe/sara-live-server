package domainsitedao

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/sys"
)

func GetDeployStateBySiteKey(ctx context.Context, siteKey string) *entity.DomainSiteDeployState {
	var row entity.DomainSiteDeployState
	if err := g.Model(string(entity.TbDomainSiteDeployState)).Ctx(ctx).
		Where("site_key = ?", strings.TrimSpace(siteKey)).
		Limit(1).
		Scan(&row); err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

// RecordDeploySuccess 原子写入站点最后成功部署时间，兼容多个服务实例同时更新。
func RecordDeploySuccess(ctx context.Context, siteKey string, uploadedAt time.Time) error {
	_, err := g.DB().Ctx(ctx).Exec(ctx, `
INSERT INTO domain_site_deploy_states (site_key, last_upload_at, created_at, updated_at)
VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE last_upload_at = VALUES(last_upload_at), updated_at = VALUES(updated_at)`,
		strings.TrimSpace(siteKey), uploadedAt, uploadedAt, uploadedAt,
	)
	return err
}
