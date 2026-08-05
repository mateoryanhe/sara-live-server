package game

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/gameplatformdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func initGameShelfCfg() {
	cfgdao.InitGameShelfCfgDao()
	cfgdao.ReloadGameShelfCfgCache()
}

// GetGameShelfList CMS 分页查询上架游戏(读内存).
func GetGameShelfList(_ context.Context, req *gameplatformdto.GameShelfListReq) (*httpserver.CMSQueryResp, error) {
	all := filterGameShelfList(cfgdao.GetAllGameShelfCfgFromMemory(), req)
	total := len(all)
	pageIndex, pageSize := normalizeCMSGameShelfPage(req.PageIndex, req.PageSize)
	start, end := cmsGameShelfPageRange(total, pageIndex, pageSize)

	list := make([]*gameplatformdto.GameShelfListItem, 0, end-start)
	for _, row := range all[start:end] {
		if row == nil {
			continue
		}
		list = append(list, toGameShelfListItem(row))
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

// AddGameShelf CMS 添加上架游戏(写库并刷新内存).
func AddGameShelf(ctx context.Context, req *gameplatformdto.AddGameShelfReq) (*gameplatformdto.AddGameShelfRes, error) {
	gameCode := strings.TrimSpace(req.GameCode)
	if gameCode == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if cfgdao.IsGameOnShelfFromMemory(gameCode) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	game, err := EnsureVendorGameForShelf(ctx, gameCode)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	row := &entity.GameShelfCfg{
		GameCode: gameCode,
	}
	row.CreatedAt = now
	row.UpdatedAt = now
	if err := cfgdao.CreateGameShelfCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadGameShelfCfgCache()
	AddOnShelfVendorGame(game)
	return &gameplatformdto.AddGameShelfRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

// DeleteGameShelf CMS 删除上架游戏(写库并刷新内存).
func DeleteGameShelf(_ context.Context, req *gameplatformdto.DeleteGameShelfReq) (*gameplatformdto.DeleteGameShelfRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	gameCode := strings.TrimSpace(req.GameCode)
	if gameCode != "" {
		if err := cfgdao.DeleteGameShelfCfgByGameCode(gameCode); err != nil {
			return nil, err
		}
		cfgdao.ReloadGameShelfCfgCache()
		RemoveOnShelfVendorGame(gameCode)
		return &gameplatformdto.DeleteGameShelfRes{Success: true}, nil
	}
	if req.ID == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	shelfRow := findGameShelfCfgByID(req.ID)
	if shelfRow != nil && shelfRow.GameCode != "" {
		gameCode = shelfRow.GameCode
	}
	if err := cfgdao.DeleteGameShelfCfg(req.ID); err != nil {
		return nil, err
	}
	cfgdao.ReloadGameShelfCfgCache()
	if gameCode != "" {
		RemoveOnShelfVendorGame(gameCode)
	}
	return &gameplatformdto.DeleteGameShelfRes{Success: true}, nil
}

// BatchAddGameShelf CMS 批量添加上架游戏(写库并刷新内存).
func BatchAddGameShelf(ctx context.Context, req *gameplatformdto.BatchAddGameShelfReq) (*gameplatformdto.BatchAddGameShelfRes, error) {
	if req == nil || len(req.GameCodes) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if err := EnsureVendorBrowseCache(ctx); err != nil {
		return nil, err
	}

	successCount := 0
	skipCount := 0
	now := time.Now()
	added := make(map[string]struct{}, len(req.GameCodes))
	for _, raw := range req.GameCodes {
		gameCode := strings.TrimSpace(raw)
		if gameCode == "" {
			continue
		}
		if cfgdao.IsGameOnShelfFromMemory(gameCode) {
			skipCount++
			continue
		}
		if _, ok := added[gameCode]; ok {
			skipCount++
			continue
		}
		game, ok := GetVendorGameFromBrowseCache(gameCode)
		if !ok {
			skipCount++
			continue
		}
		row := &entity.GameShelfCfg{GameCode: gameCode}
		row.CreatedAt = now
		row.UpdatedAt = now
		if err := cfgdao.CreateGameShelfCfg(row); err != nil {
			return nil, err
		}
		added[gameCode] = struct{}{}
		AddOnShelfVendorGame(game)
		successCount++
	}
	cfgdao.ReloadGameShelfCfgCache()
	return &gameplatformdto.BatchAddGameShelfRes{
		Success:      true,
		SuccessCount: successCount,
		SkipCount:    skipCount,
	}, nil
}

// BatchDeleteGameShelf CMS 批量删除上架游戏(写库并刷新内存).
func BatchDeleteGameShelf(_ context.Context, req *gameplatformdto.BatchDeleteGameShelfReq) (*gameplatformdto.BatchDeleteGameShelfRes, error) {
	if req == nil || len(req.GameCodes) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	codes := make([]string, 0, len(req.GameCodes))
	for _, raw := range req.GameCodes {
		code := strings.TrimSpace(raw)
		if code != "" {
			codes = append(codes, code)
		}
	}
	count, err := cfgdao.DeleteGameShelfCfgByGameCodes(codes)
	if err != nil {
		return nil, err
	}
	cfgdao.ReloadGameShelfCfgCache()
	RemoveOnShelfVendorGames(codes)
	return &gameplatformdto.BatchDeleteGameShelfRes{
		Success:      true,
		SuccessCount: count,
	}, nil
}

// GetAllOnShelfVendorGamesFromMemory 获取已上架的第三方游戏(内存).
func GetAllOnShelfVendorGamesFromMemory() []*VendorGame {
	return GetAllVendorGamesFromMemory()
}

func findGameShelfCfgByID(id uint64) *entity.GameShelfCfg {
	for _, row := range cfgdao.GetAllGameShelfCfgFromMemory() {
		if row != nil && row.ID == id {
			return row
		}
	}
	return nil
}

func filterGameShelfList(all []*entity.GameShelfCfg, req *gameplatformdto.GameShelfListReq) []*entity.GameShelfCfg {
	if req == nil {
		return all
	}
	gameCode := strings.TrimSpace(req.GameCode)
	if gameCode == "" {
		return all
	}
	list := make([]*entity.GameShelfCfg, 0, len(all))
	for _, row := range all {
		if row == nil {
			continue
		}
		if !strings.Contains(strings.ToLower(row.GameCode), strings.ToLower(gameCode)) {
			continue
		}
		list = append(list, row)
	}
	return list
}

func toGameShelfListItem(row *entity.GameShelfCfg) *gameplatformdto.GameShelfListItem {
	return &gameplatformdto.GameShelfListItem{
		ID:       strconv.FormatUint(row.ID, 10),
		GameCode: row.GameCode,
	}
}

const (
	cmsGameShelfDefaultPageSize = 20
	cmsGameShelfMaxPageSize     = 100
)

func normalizeCMSGameShelfPage(pageIndex, pageSize int) (int, int) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = cmsGameShelfDefaultPageSize
	}
	if pageSize > cmsGameShelfMaxPageSize {
		pageSize = cmsGameShelfMaxPageSize
	}
	return pageIndex, pageSize
}

func cmsGameShelfPageRange(total, pageIndex, pageSize int) (int, int) {
	start := (pageIndex - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	return start, end
}
