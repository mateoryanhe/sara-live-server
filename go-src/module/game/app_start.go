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
	platform := strings.TrimSpace(req.Platform)
	if gameCode == "" || platform == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if !cfgdao.IsGameOnShelfFromMemory(gameCode) {
		return nil, errercode.CreateCode(errercode.GameCfgNonExist)
	}

	lang := strings.TrimSpace(req.Lang)
	link, err := fetchVendorGameStartURL(ctx, gameCode, platform, strconv.FormatUint(userId, 10), lang)
	if err != nil {
		xrlog.ErrorWithErr(ctx, "Game", "fetch vendor game start url failed", err)
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	return &gameplatformdto.AppGameStartRes{Link: link}, nil
}
