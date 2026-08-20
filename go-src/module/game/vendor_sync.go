package game

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/dao/cfgdao"
	entity "xr-game-server/entity/game"
	"xr-game-server/errercode"
)

// SyncVendorGameLibraryFromVendor 从第三方全量拉取并覆盖游戏库表.
func SyncVendorGameLibraryFromVendor(ctx context.Context) (int, error) {
	if ctx == nil {
		ctx = gctx.New()
	}
	games, err := fetchAllVendorGames(ctx)
	if err != nil {
		return 0, err
	}
	now := time.Now()
	rows := make([]*entity.VendorGameLib, 0, len(games))
	for _, game := range games {
		if row := toVendorGameLibEntity(game); row != nil {
			rows = append(rows, row)
		}
	}
	rows = cfgdao.BuildVendorGameLibRows(rows, now)
	if err := cfgdao.ReplaceAllVendorGameLibs(rows); err != nil {
		return 0, err
	}
	repairShelfPlatformFromVendorLibrary()
	vendorDetailLog().Infof(ctx, "sync vendor game library done total=%d", len(rows))
	return len(rows), nil
}

// GetVendorGameFromLibrary 从游戏库表读取游戏.
func GetVendorGameFromLibrary(gameCode, platform string) (*VendorGame, bool) {
	gameCode = strings.TrimSpace(gameCode)
	platform = strings.TrimSpace(platform)
	if gameCode == "" || platform == "" {
		return nil, false
	}
	row := cfgdao.GetVendorGameLib(gameCode, platform)
	if row == nil {
		return nil, false
	}
	return toVendorGame(row), true
}

// EnsureVendorGameForShelf 确认游戏库中存在该游戏，用于上架.
func EnsureVendorGameForShelf(_ context.Context, gameCode, platform string) (*VendorGame, error) {
	gameCode = strings.TrimSpace(gameCode)
	platform = strings.TrimSpace(platform)
	if gameCode == "" || platform == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	game, ok := GetVendorGameFromLibrary(gameCode, platform)
	if !ok || game == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return game, nil
}
