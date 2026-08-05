package game

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/xrlog"
	"xr-game-server/core/xrpool"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/errercode"
)

// ReloadVendorGameCache 启动/保存平台配置时拉取全量第三方数据并同步上架缓存.
func ReloadVendorGameCache() {
	reloadVendorGameCacheWithRetry(gctx.New())
}

// ReloadVendorGameCacheAsync 异步刷新第三方浏览缓存(失败自动重试).
func ReloadVendorGameCacheAsync(ctx context.Context) {
	if ctx == nil {
		ctx = gctx.New()
	}
	xrpool.AddWithRecover(ctx, func(runCtx context.Context) {
		reloadVendorGameCacheWithRetry(runCtx)
	})
}

func reloadVendorGameCacheWithRetry(ctx context.Context) {
	if ctx == nil {
		ctx = gctx.New()
	}
	var lastErr error
	for attempt := 1; attempt <= vendorGameListMaxAttempts; attempt++ {
		if attempt > 1 {
			vendorDetailLog().Warningf(ctx, "reload vendor game cache retry attempt=%d/%d wait=%s err=%v",
				attempt, vendorGameListMaxAttempts, vendorGameListRetryInterval, lastErr)
			if err := sleepWithContext(ctx, vendorGameListRetryInterval); err != nil {
				return
			}
		}
		if err := refreshVendorBrowseCache(ctx, true); err == nil {
			return
		} else {
			lastErr = err
		}
	}
	vendorDetailLog().Warningf(ctx, "reload vendor game cache gave up after %d attempts err=%v",
		vendorGameListMaxAttempts, lastErr)
	xrlog.ErrorWithErr(ctx, "VendorGame", "reload vendor game cache failed", lastErr)
}

// EnsureVendorBrowseCache 浏览缓存未命中或过期时拉取全量第三方数据(30 分钟有效).
func EnsureVendorBrowseCache(ctx context.Context) error {
	if !isVendorBrowseCacheEmpty() {
		return nil
	}
	err := refreshVendorBrowseCache(ctx, false)
	if err != nil {
		xrlog.ErrorWithErr(ctx, "VendorGame", "refresh vendor browse cache failed", err)
	}
	return err
}

func refreshVendorBrowseCache(ctx context.Context, syncOnShelf bool) error {
	if ctx == nil {
		ctx = gctx.New()
	}
	start := time.Now()
	vendorDetailLog().Infof(ctx, "refresh vendor browse cache start syncOnShelf=%v", syncOnShelf)

	games, err := fetchAllVendorGames(ctx)
	totalCostMs := time.Since(start).Milliseconds()
	if err != nil {
		vendorDetailLog().Warningf(ctx, "refresh vendor browse cache failed totalCostMs=%d err=%v", totalCostMs, err)
		return err
	}

	setVendorBrowseCache(games)
	vendorMap := buildVendorGameMap(games)

	removedCount := 0
	if syncOnShelf {
		shelfSet := cfgdao.GetGameShelfCodeSetFromMemory()
		if len(shelfSet) > 0 {
			removedCount = pruneInvalidShelfGames(ctx, shelfSet, vendorMap)
		}
		refreshOnShelfVendorGamesFromMap(vendorMap)
	} else {
		refreshOnShelfVendorDetailsFromMap(vendorMap)
	}

	vendorDetailLog().Infof(ctx,
		"refresh vendor browse cache done vendorTotal=%d onShelf=%d removedShelf=%d totalCostMs=%d",
		len(games), len(getOnShelfVendorGameCodesSnapshot()), removedCount, totalCostMs)
	logVendorGameListData(ctx, games)
	return nil
}

func syncVendorBrowseCache(syncOnShelf bool) {
	_ = refreshVendorBrowseCache(gctx.New(), syncOnShelf)
}

// EnsureVendorGameForShelf 确认第三方存在该游戏，用于上架.
func EnsureVendorGameForShelf(ctx context.Context, gameCode string) (*VendorGame, error) {
	gameCode = strings.TrimSpace(gameCode)
	if gameCode == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := EnsureVendorBrowseCache(ctx); err != nil {
		return nil, err
	}
	if game, ok := GetVendorGameFromBrowseCache(gameCode); ok {
		return game, nil
	}
	return nil, errercode.CreateCode(errercode.InvalidParam)
}

func refreshOnShelfVendorDetailsFromMap(vendorMap map[string]*VendorGame) {
	shelfSet := cfgdao.GetGameShelfCodeSetFromMemory()
	if len(shelfSet) == 0 {
		return
	}
	for code := range shelfSet {
		if game, ok := vendorMap[code]; ok {
			AddOnShelfVendorGame(game)
		}
	}
}

func buildVendorGameMap(games []*VendorGame) map[string]*VendorGame {
	m := make(map[string]*VendorGame, len(games))
	for _, row := range games {
		if row == nil {
			continue
		}
		code := strings.TrimSpace(row.GameCode)
		if code == "" {
			continue
		}
		m[code] = row
	}
	return m
}

func pruneInvalidShelfGames(ctx context.Context, shelfSet map[string]struct{}, vendorMap map[string]*VendorGame) int {
	if len(shelfSet) == 0 {
		return 0
	}
	removed := make([]string, 0)
	for code := range shelfSet {
		if _, ok := vendorMap[code]; ok {
			continue
		}
		if err := cfgdao.DeleteGameShelfCfgByGameCode(code); err != nil {
			vendorDetailLog().Warningf(ctx, "prune invalid shelf game failed gameCode=%s err=%v", code, err)
			continue
		}
		RemoveOnShelfVendorGame(code)
		removed = append(removed, code)
	}
	if len(removed) == 0 {
		return 0
	}
	cfgdao.ReloadGameShelfCfgCache()
	vendorDetailLog().Infof(ctx, "pruned invalid shelf games count=%d codes=%v", len(removed), removed)
	return len(removed)
}

func logVendorGameListData(ctx context.Context, games []*VendorGame) {
	raw, err := json.Marshal(games)
	if err != nil {
		vendorDetailLog().Warningf(ctx, "vendor game list data marshal failed: %v", err)
		return
	}
	vendorDetailLog().Infof(ctx, "vendor game list data: %s", string(raw))
}
