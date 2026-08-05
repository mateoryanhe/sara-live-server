package game

import (
	"sort"
	"sync"

	"xr-game-server/core/cache"
	"xr-game-server/dao/cfgdao"
)

const vendorBrowseCacheKey = "all"

// VendorGame 第三方平台游戏
type VendorGame struct {
	GameCode string `json:"gameCode"`
	Name     string `json:"name"`
	NameEn   string `json:"nameEn"`
	Category string `json:"category"`
	Cover    string `json:"cover"`
	Platform string `json:"platform"`
}

var (
	vendorGameCacheMu    sync.RWMutex
	vendorBrowseCacheMgr *cache.CacheMgr
	// vendorOnShelfCache 已上架游戏详情，永不过期.
	vendorOnShelfCache = make(map[string]*VendorGame)
)

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

func isVendorBrowseCacheEmpty() bool {
	return len(GetAllVendorBrowseGamesFromMemory()) == 0
}

func AddOnShelfVendorGame(game *VendorGame) {
	if game == nil || game.GameCode == "" {
		return
	}
	vendorGameCacheMu.Lock()
	defer vendorGameCacheMu.Unlock()
	vendorOnShelfCache[game.GameCode] = cloneVendorGame(game)
}

func RemoveOnShelfVendorGame(gameCode string) {
	if gameCode == "" {
		return
	}
	vendorGameCacheMu.Lock()
	defer vendorGameCacheMu.Unlock()
	delete(vendorOnShelfCache, gameCode)
}

func RemoveOnShelfVendorGames(gameCodes []string) {
	if len(gameCodes) == 0 {
		return
	}
	vendorGameCacheMu.Lock()
	defer vendorGameCacheMu.Unlock()
	for _, code := range gameCodes {
		if code == "" {
			continue
		}
		delete(vendorOnShelfCache, code)
	}
}

func refreshOnShelfVendorGamesFromMap(vendorMap map[string]*VendorGame) {
	shelfRows := cfgdao.GetAllGameShelfCfgFromMemory()
	vendorGameCacheMu.Lock()
	defer vendorGameCacheMu.Unlock()
	next := make(map[string]*VendorGame, len(shelfRows))
	for _, row := range shelfRows {
		if row == nil || row.GameCode == "" {
			continue
		}
		if game, ok := vendorMap[row.GameCode]; ok {
			next[row.GameCode] = cloneVendorGame(game)
		}
	}
	vendorOnShelfCache = next
}

// GetAllVendorGamesFromMemory 获取已上架第三方游戏(内存快照，永不过期).
func GetAllVendorGamesFromMemory() []*VendorGame {
	shelfRows := cfgdao.GetAllGameShelfCfgFromMemory()
	vendorGameCacheMu.RLock()
	defer vendorGameCacheMu.RUnlock()
	list := make([]*VendorGame, 0, len(shelfRows))
	for _, row := range shelfRows {
		if row == nil || row.GameCode == "" {
			continue
		}
		if game, ok := vendorOnShelfCache[row.GameCode]; ok {
			list = append(list, cloneVendorGame(game))
		}
	}
	return list
}

func getOnShelfVendorGameCodesSnapshot() []string {
	vendorGameCacheMu.RLock()
	defer vendorGameCacheMu.RUnlock()
	codes := make([]string, 0, len(vendorOnShelfCache))
	for code := range vendorOnShelfCache {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	return codes
}
