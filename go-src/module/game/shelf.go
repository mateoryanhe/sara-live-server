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

// GetGameShelfList CMS 分页查询上架游戏(读内存).
func GetGameShelfList(_ context.Context, req *gameplatformdto.GameShelfListReq) (*httpserver.CMSQueryResp, error) {
	all := filterGameShelfList(cfgdao.GetAllGameCfgFromMemory(), req)
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

// AddGameShelf CMS 上架游戏: 写 game_cfgs 并刷新永久缓存.
func AddGameShelf(ctx context.Context, req *gameplatformdto.AddGameShelfReq) (*gameplatformdto.AddGameShelfRes, error) {
	gameCode := strings.TrimSpace(req.GameCode)
	if gameCode == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if cfgdao.IsGameOnShelfFromMemory(gameCode) {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}

	vendorGame, err := EnsureVendorGameForShelf(ctx, gameCode)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	row := newShelfGameCfg(vendorGame, gameCode, now)
	if err := cfgdao.CreateGameCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadGameCfgCache()
	return &gameplatformdto.AddGameShelfRes{
		Success: true,
		ID:      strconv.FormatUint(row.ID, 10),
	}, nil
}

// DeleteGameShelf CMS 下架游戏: 删 game_cfgs 并刷新永久缓存.
func DeleteGameShelf(_ context.Context, req *gameplatformdto.DeleteGameShelfReq) (*gameplatformdto.DeleteGameShelfRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	gameCode := strings.TrimSpace(req.GameCode)
	if gameCode != "" {
		if err := cfgdao.DeleteGameCfgByGameCode(gameCode); err != nil {
			return nil, err
		}
		cfgdao.ReloadGameCfgCache()
		return &gameplatformdto.DeleteGameShelfRes{Success: true}, nil
	}
	if req.ID == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	shelfRow := findGameCfgByID(req.ID)
	if shelfRow != nil && shelfRow.GameCode != "" {
		gameCode = shelfRow.GameCode
	}
	if err := cfgdao.DeleteGameCfg(req.ID); err != nil {
		return nil, err
	}
	cfgdao.ReloadGameCfgCache()
	return &gameplatformdto.DeleteGameShelfRes{Success: true}, nil
}

// BatchAddGameShelf CMS 批量上架(写 game_cfgs 并刷新永久缓存).
func BatchAddGameShelf(_ context.Context, req *gameplatformdto.BatchAddGameShelfReq) (*gameplatformdto.BatchAddGameShelfRes, error) {
	if req == nil || len(req.GameCodes) == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
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
		vendorGame, ok := GetVendorGameFromBrowseCache(gameCode)
		if !ok {
			skipCount++
			continue
		}
		row := newShelfGameCfg(vendorGame, gameCode, now)
		if err := cfgdao.CreateGameCfg(row); err != nil {
			return nil, err
		}
		added[gameCode] = struct{}{}
		successCount++
	}
	cfgdao.ReloadGameCfgCache()
	return &gameplatformdto.BatchAddGameShelfRes{
		Success:      true,
		SuccessCount: successCount,
		SkipCount:    skipCount,
	}, nil
}

// BatchDeleteGameShelf CMS 批量下架(删 game_cfgs 并刷新永久缓存).
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
	count, err := cfgdao.DeleteGameCfgByGameCodes(codes)
	if err != nil {
		return nil, err
	}
	cfgdao.ReloadGameCfgCache()
	return &gameplatformdto.BatchDeleteGameShelfRes{
		Success:      true,
		SuccessCount: count,
	}, nil
}

func findGameCfgByID(id uint64) *entity.GameCfg {
	for _, row := range cfgdao.GetAllGameCfgFromMemory() {
		if row != nil && row.ID == id {
			return row
		}
	}
	return nil
}

func filterGameShelfList(all []*entity.GameCfg, req *gameplatformdto.GameShelfListReq) []*entity.GameCfg {
	if req == nil {
		return all
	}
	gameCode := strings.TrimSpace(req.GameCode)
	if gameCode == "" {
		return all
	}
	list := make([]*entity.GameCfg, 0, len(all))
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

func toGameShelfListItem(row *entity.GameCfg) *gameplatformdto.GameShelfListItem {
	return &gameplatformdto.GameShelfListItem{
		ID:       strconv.FormatUint(row.ID, 10),
		GameCode: row.GameCode,
		Platform: row.Platform,
	}
}

func newShelfGameCfg(vendorGame *VendorGame, gameCode string, now time.Time) *entity.GameCfg {
	row := &entity.GameCfg{
		GameCode: gameCode,
		Cover:    strings.TrimSpace(vendorGame.Cover),
		NameEn:   strings.TrimSpace(vendorGame.NameEn),
		Platform: strings.TrimSpace(vendorGame.Platform),
	}
	row.CreatedAt = now
	row.UpdatedAt = now
	return row
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
