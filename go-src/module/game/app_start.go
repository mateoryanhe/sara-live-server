package game

import (
	"context"
	"strconv"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/core/xrlog"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/gameplatformdto"
	"xr-game-server/errercode"
)

// GetAppGameStartLink App 获取游戏启动链接(调用第三方 /open/game/start).
func GetAppGameStartLink(ctx context.Context, req *gameplatformdto.AppGameStartReq) (*gameplatformdto.AppGameStartRes, error) {
	userId := httpserver.GetAuthId(ctx)
	if userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	gameCode := strings.TrimSpace(req.GameCode)
	if gameCode == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if !cfgdao.IsGameOnShelfFromMemory(gameCode) {
		return nil, errercode.CreateCode(errercode.GameCfgNonExist)
	}

	platform, err := resolveGameStartPlatform(gameCode)
	if err != nil {
		return nil, err
	}

	link, err := fetchVendorGameStartURL(ctx, gameCode, platform, strconv.FormatUint(userId, 10), "en")
	if err != nil {
		xrlog.ErrorWithErr(ctx, "Game", "fetch vendor game start url failed", err)
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return &gameplatformdto.AppGameStartRes{Link: link}, nil
}

func resolveGameStartPlatform(gameCode string) (string, error) {
	gameCfg := cfgdao.GetGameCfgByGameCode(gameCode)
	if gameCfg == nil {
		return "", errercode.CreateCode(errercode.GameCfgNonExist)
	}
	platform := strings.TrimSpace(gameCfg.Platform)
	if platform == "" {
		return "", errercode.CreateCode(errercode.InvalidParam)
	}
	return platform, nil
}
