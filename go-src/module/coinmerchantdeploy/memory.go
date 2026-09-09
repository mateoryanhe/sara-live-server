package coinmerchantdeploy

import (
	"strings"
	"sync/atomic"
	"time"

	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity/sys"
)

type cfgSnapshot struct {
	ID           uint64
	DeploySecret string
	UpdatedAt    string
}

var cfgCache atomic.Value

func reloadCfgMemory() {
	cfgCache.Store(toCfgSnapshot(cfgdao.LoadCoinMerchantDeployCfg()))
}

func ReloadCoinMerchantDeployCache() {
	reloadCfgMemory()
}

func getCfgCache() *cfgSnapshot {
	if cfgCache.Load() == nil {
		reloadCfgMemory()
	}
	v := cfgCache.Load()
	if v == nil {
		return &cfgSnapshot{}
	}
	snap, ok := v.(*cfgSnapshot)
	if !ok || snap == nil {
		return &cfgSnapshot{}
	}
	return snap
}

// GetDeploySecret 币商 H5 加解密密钥(Header X-Coin-Merchant-Client=1 时使用)
func GetDeploySecret() string {
	return getCfgCache().DeploySecret
}

func toCfgSnapshot(row *entity.CoinMerchantDeployCfg) *cfgSnapshot {
	if row == nil || row.ID == 0 {
		return &cfgSnapshot{}
	}
	return &cfgSnapshot{
		ID:           row.ID,
		DeploySecret: strings.TrimSpace(row.DeploySecret),
		UpdatedAt:    formatCfgTime(row.UpdatedAt),
	}
}

func formatCfgTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
