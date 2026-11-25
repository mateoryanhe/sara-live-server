package shareappdownload

import (
	"strings"
	"sync/atomic"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/user"
)

type sofieShareRowSnapshot struct {
	PackageName   string
	GooglePlayUrl string
	ShareText     string
}

var (
	sofieShareCfgCache         atomic.Value // *sofieShareRowSnapshot
	emptySofieShareCfgSnapshot = &sofieShareRowSnapshot{}
)

func reloadSofieShareMemory() {
	sofieShareCfgCache.Store(toSofieShareRowSnapshot(cfgdao.LoadSofieShareAppDownloadCfg()))
}

func getSofieShareCfg() *sofieShareRowSnapshot {
	v := sofieShareCfgCache.Load()
	if v == nil {
		return emptySofieShareCfgSnapshot
	}
	cfg, ok := v.(*sofieShareRowSnapshot)
	if !ok || cfg == nil {
		return emptySofieShareCfgSnapshot
	}
	return cfg
}

func toSofieShareRowSnapshot(row *entity.SofieShareAppDownloadCfg) *sofieShareRowSnapshot {
	if row == nil || row.ID == 0 {
		return emptySofieShareCfgSnapshot
	}
	return &sofieShareRowSnapshot{
		PackageName:   strings.TrimSpace(row.PackageName),
		GooglePlayUrl: strings.TrimSpace(row.GooglePlayUrl),
		ShareText:     strings.TrimSpace(row.ShareText),
	}
}

// ResolveSofieShareSnapshot 单行配置: 包名一致或任意未命中均回落该行(当前仅一个 Sofie 包).
func ResolveSofieShareSnapshot(packageName string) *sofieShareRowSnapshot {
	cfg := getSofieShareCfg()
	if cfg.GooglePlayUrl == "" && cfg.ShareText == "" {
		return nil
	}
	_ = strings.TrimSpace(packageName)
	return cfg
}

// ResolveSofieGooglePlayUrl 解析 Google Play 跳转目标.
func ResolveSofieGooglePlayUrl(packageName string) string {
	snap := ResolveSofieShareSnapshot(packageName)
	if snap == nil {
		return ""
	}
	return strings.TrimSpace(snap.GooglePlayUrl)
}