package game

import (
	"strings"

	"xr-game-server/dao/cfgdao"
)

// repairShelfMetadataFromVendorLibrary 游戏库同步后更新已上架游戏的第三方元数据.
func repairShelfMetadataFromVendorLibrary() {
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
		cover := strings.TrimSpace(libRow.Cover)
		if platform != "" && platform != strings.TrimSpace(row.Platform) {
			ok, err := cfgdao.SetGameCfgPlatform(row.GameCode, platform)
			if err == nil && ok {
				updated = true
			}
		}
		if cover != strings.TrimSpace(row.Cover) {
			ok, err := cfgdao.SetGameCfgCover(row.GameCode, cover)
			if err == nil && ok {
				updated = true
			}
		}
	}
	if updated {
		cfgdao.ReloadGameCfgCache()
	}
}
