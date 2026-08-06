package game

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/errercode"
)

// ForceRefreshVendorBrowseCache CMS 查询时从第三方全量拉取并覆盖 30 分钟浏览缓存.
func ForceRefreshVendorBrowseCache(ctx context.Context) error {
	if ctx == nil {
		ctx = gctx.New()
	}
	games, err := fetchAllVendorGames(ctx)
	if err != nil {
		return err
	}
	setVendorBrowseCache(games)
	vendorDetailLog().Infof(ctx, "refresh vendor browse cache done vendorTotal=%d", len(games))
	return nil
}

// EnsureVendorGameForShelf 确认第三方浏览缓存中存在该游戏，用于上架.
func EnsureVendorGameForShelf(_ context.Context, gameCode string) (*VendorGame, error) {
	gameCode = strings.TrimSpace(gameCode)
	if gameCode == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	game, ok := GetVendorGameFromBrowseCache(gameCode)
	if !ok || game == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return game, nil
}
