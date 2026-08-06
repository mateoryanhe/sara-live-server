package game

import (
	"xr-game-server/core/cache"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/entity"
)

const vendorBrowseCacheKey = "all"

// VendorGame 第三方平台游戏(仅 CMS 浏览缓存使用)
type VendorGame struct {
	GameCode string `json:"gameCode"`
	Name     string `json:"name"`
	NameEn   string `json:"nameEn"`
	Category string `json:"category"`
	Cover    string `json:"cover"`
	Platform string `json:"platform"`
}

var vendorBrowseCacheMgr *cache.CacheMgr

func initVendorGameCache() {
	vendorBrowseCacheMgr = cache.NewCacheMgr()
}

func cloneVendorGame(game *VendorGame) *VendorGame {
	if game == nil {
		return nil
	}
	return &VendorGame{
		GameCode: game.GameCode,
		Name:     game.Name,
		NameEn:   game.NameEn,
		Category: game.Category,
		Cover:    game.Cover,
		Platform: game.Platform,
	}
}

func setVendorBrowseCache(games []*VendorGame) {
	if vendorBrowseCacheMgr == nil {
		return
	}
	snapshot := make([]*VendorGame, 0, len(games))
	for _, row := range games {
		if row == nil {
			continue
		}
		snapshot = append(snapshot, cloneVendorGame(row))
	}
	vendorBrowseCacheMgr.FlushCache(vendorBrowseCacheKey, snapshot)
}

func GetAllVendorBrowseGamesFromMemory() []*VendorGame {
	if vendorBrowseCacheMgr == nil {
		return make([]*VendorGame, 0)
	}
	v := vendorBrowseCacheMgr.GetFromCache(vendorBrowseCacheKey)
	if v == nil {
		return make([]*VendorGame, 0)
	}
	list, _ := v.([]*VendorGame)
	if len(list) == 0 {
		return make([]*VendorGame, 0)
	}
	out := make([]*VendorGame, 0, len(list))
	for _, row := range list {
		out = append(out, cloneVendorGame(row))
	}
	return out
}

func GetVendorGameFromBrowseCache(gameCode string) (*VendorGame, bool) {
	for _, row := range GetAllVendorBrowseGamesFromMemory() {
		if row != nil && row.GameCode == gameCode {
			return cloneVendorGame(row), true
		}
	}
	return nil, false
}

// GetAllOnShelfGamesFromMemory 获取已上架游戏(读 game_cfgs 永久缓存).
func GetAllOnShelfGamesFromMemory() []*entity.GameCfg {
	return cfgdao.GetAllGameCfgFromMemory()
}
