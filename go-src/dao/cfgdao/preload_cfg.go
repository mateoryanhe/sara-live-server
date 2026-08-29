package cfgdao

import (
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/entity/sys"
)

const (
	DefaultRecentLoginPreloadLimit = 100
	DefaultHotRestartAuth          = "nGH66S4TjBjQqCKyWJAM"
	DefaultMemoryLimitM            = 300
	DefaultIpGeoDbPath             = "/home/ec2-user/xgameserver/GeoLite2-Country.mmdb"
	DefaultInitGold                = 0
	DefaultInitDiamond             = 0
)

func LoadPreloadCfg() *entity.PreloadCfg {
	var row entity.PreloadCfg
	if err := g.DB().Model(string(entity.TbPreloadCfg)).Order("id asc").Limit(1).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func SavePreloadCfg(row *entity.PreloadCfg) error {
	if row == nil {
		return nil
	}
	_, err := g.DB().Model(string(entity.TbPreloadCfg)).Save(row)
	return err
}

// GetRecentLoginPreloadLimit 读取最近登录用户预热数量,未配置或无效时返回默认值
func GetRecentLoginPreloadLimit() int {
	cfg := LoadPreloadCfg()
	if cfg == nil || cfg.RecentLoginLimit <= 0 {
		return DefaultRecentLoginPreloadLimit
	}
	return cfg.RecentLoginLimit
}

func GetHotRestartAuth() string {
	cfg := LoadPreloadCfg()
	if cfg != nil {
		if auth := strings.TrimSpace(cfg.HotRestartAuth); auth != "" {
			return auth
		}
	}
	return DefaultHotRestartAuth
}

func GetMemoryLimitM() int {
	cfg := LoadPreloadCfg()
	if cfg != nil && cfg.MemoryLimitM > 0 {
		return cfg.MemoryLimitM
	}
	return DefaultMemoryLimitM
}

func GetIpGeoDbPath() string {
	cfg := LoadPreloadCfg()
	if cfg != nil {
		if path := strings.TrimSpace(cfg.IpGeoDbPath); path != "" {
			return path
		}
	}
	return DefaultIpGeoDbPath
}

func GetInitGold() float64 {
	cfg := LoadPreloadCfg()
	if cfg == nil || cfg.InitGold < 0 {
		return DefaultInitGold
	}
	return cfg.InitGold
}

func GetInitDiamond() float64 {
	cfg := LoadPreloadCfg()
	if cfg == nil || cfg.InitDiamond < 0 {
		return DefaultInitDiamond
	}
	return cfg.InitDiamond
}
