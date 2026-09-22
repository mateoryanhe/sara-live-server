package domainsite

import (
	"context"
	"strings"
	"time"

	"xr-game-server/dao/domainsitedao"
)

const deployTimeLayout = "2006-01-02 15:04:05"

// RecordDeploySuccess 仅在站点文件完整写入并解压成功后调用。
func RecordDeploySuccess(ctx context.Context, siteKey string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	siteKey = strings.TrimSpace(siteKey)
	if siteKey == "" {
		return "", nil
	}
	uploadedAt := time.Now()
	if err := domainsitedao.RecordDeploySuccess(ctx, siteKey, uploadedAt); err != nil {
		return "", err
	}
	return uploadedAt.Format(deployTimeLayout), nil
}

func GetLastUploadAt(ctx context.Context, siteKey string) string {
	if ctx == nil {
		ctx = context.Background()
	}
	row := domainsitedao.GetDeployStateBySiteKey(ctx, siteKey)
	if row == nil || row.LastUploadAt.IsZero() {
		return ""
	}
	return row.LastUploadAt.Format(deployTimeLayout)
}
