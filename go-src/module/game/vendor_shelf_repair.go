package game

import (
	"strings"

	"xr-game-server/dao/cfgdao"
)

// repairShelfPlatformFromVendorLibrary 游戏库同步后, 用 vendorGame.Platform 修正已上架记录.
func repairShelfPlatformFromVendorLibrary() {
	updated := false
	for _, row := range cfgdao.GetAllGameCfgFromMemory() {
		if row == nil {
			continue
		}
		libRow := cfgdao.GetVendorGameLib(row.GameCode, row.Platform)
		if libRow == nil {
			continue
		}
		platform := strings.TrimSpace(libRow.Platform)
		if platform == "" || platform == strings.TrimSpace(row.Platform) {
			continue
		}
		ok, err := cfgdao.SetGameCfgPlatform(row.GameCode, platform)
		if err != nil || !ok {
			continue
		}
		updated = true
	}
	if updated {
		cfgdao.ReloadGameCfgCache()
	}
}
