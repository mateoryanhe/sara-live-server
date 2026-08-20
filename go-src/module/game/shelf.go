package game

import (
	"context"
	"strconv"
	"strings"
	"time"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/accountdao"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/gameplatformdto"
	"xr-game-server/entity/game"
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

// UpdateGameShelf CMS 更新上架游戏直播展示字段.
func UpdateGameShelf(_ context.Context, req *gameplatformdto.UpdateGameShelfReq) (*gameplatformdto.UpdateGameShelfRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	gameCode := strings.TrimSpace(req.GameCode)
	if gameCode == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if cfgdao.GetGameCfgByGameCode(gameCode) == nil {
		return nil, errercode.CreateCode(errercode.GameCfgNonExist)
	}
	liveGameName := strings.TrimSpace(req.LiveGameName)
	liveGameCover := strings.TrimSpace(req.LiveGameCover)
	if liveGameCover != "" && !strings.HasPrefix(liveGameCover, "http://") && !strings.HasPrefix(liveGameCover, "https://") {
		liveGameCover = normalizeVendorGameCover(liveGameCover)
	}
	ok, err := cfgdao.UpdateGameCfgLiveDisplay(gameCode, liveGameName, liveGameCover)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errercode.CreateCode(errercode.GameCfgNonExist)
	}
	cfgdao.ReloadGameCfgCache()
	return &gameplatformdto.UpdateGameShelfRes{Success: true}, nil
}

// GetMultiplayerConfigUrl CMS 获取第三方自研游戏 embed 配置页链接.
func GetMultiplayerConfigUrl(ctx context.Context, req *gameplatformdto.GetMultiplayerConfigUrlReq) (*gameplatformdto.GetMultiplayerConfigUrlRes, error) {
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

	platform := strings.TrimSpace(req.Platform)
	if platform == "" {
		if row := cfgdao.GetGameCfgByGameCode(gameCode); row != nil {
			platform = strings.TrimSpace(row.Platform)
		}
	}
	if platform == "" {
		platform = vendorMultiplayerDefaultPlatform
	}

	configURL, expireInMs, err := fetchVendorMultiplayerConfigURL(ctx, gameCode, platform)
	if err != nil {
		return nil, err
	}
	return &gameplatformdto.GetMultiplayerConfigUrlRes{
		ConfigUrl:  configURL,
		ExpireInMs: expireInMs,
	}, nil
}

// GetCMSGameStartLink CMS 代指定用户获取游戏启动链接.
func GetCMSGameStartLink(ctx context.Context, req *gameplatformdto.CMSGameStartLinkReq) (*gameplatformdto.CMSGameStartLinkRes, error) {
	if req == nil || req.UserId == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	gameCode := strings.TrimSpace(req.GameCode)
	if gameCode == "" {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if accountdao.GetAccountById(req.UserId) == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if !cfgdao.IsGameOnShelfFromMemory(gameCode) {
		return nil, errercode.CreateCode(errercode.GameCfgNonExist)
	}

	platform, err := resolveGameStartPlatform(gameCode)
	if err != nil {
		return nil, err
	}

	link, err := fetchVendorGameStartURL(ctx, gameCode, platform, strconv.FormatUint(req.UserId, 10), "en")
	if err != nil {
		return nil, err
	}
	return &gameplatformdto.CMSGameStartLinkRes{Link: link}, nil
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
	name := strings.TrimSpace(req.Name)
	platform := strings.TrimSpace(req.Platform)
	if gameCode == "" && name == "" && platform == "" {
		return all
	}

	list := make([]*entity.GameCfg, 0, len(all))
	for _, row := range all {
		if row == nil {
			continue
		}
		if gameCode != "" && !strings.Contains(strings.ToLower(row.GameCode), strings.ToLower(gameCode)) {
			continue
		}
		if platform != "" && !strings.Contains(strings.ToLower(row.Platform), strings.ToLower(platform)) {
			continue
		}
		if name != "" {
			nameEn := strings.ToLower(strings.TrimSpace(row.NameEn))
			vendorName := ""
			vendorNameEn := ""
			if vendorGame, ok := GetVendorGameFromBrowseCache(row.GameCode); ok && vendorGame != nil {
				vendorName = strings.ToLower(strings.TrimSpace(vendorGame.Name))
				vendorNameEn = strings.ToLower(strings.TrimSpace(vendorGame.NameEn))
			}
			keyword := strings.ToLower(name)
			if !strings.Contains(nameEn, keyword) &&
				!strings.Contains(vendorName, keyword) &&
				!strings.Contains(vendorNameEn, keyword) {
				continue
			}
		}
		list = append(list, row)
	}
	return list
}

func toGameShelfListItem(row *entity.GameCfg) *gameplatformdto.GameShelfListItem {
	item := &gameplatformdto.GameShelfListItem{
		ID:               strconv.FormatUint(row.ID, 10),
		GameCode:         row.GameCode,
		NameEn:           row.NameEn,
		Cover:            BuildGameCoverUrl(row.Cover),
		LiveGameName:     row.LiveGameName,
		LiveGameCover:    row.LiveGameCover,
		LiveGameCoverUrl: BuildGameCoverUrl(row.LiveGameCover),
		Platform:         row.Platform,
	}
	if vendorGame, ok := GetVendorGameFromBrowseCache(row.GameCode); ok && vendorGame != nil {
		item.Name = strings.TrimSpace(vendorGame.Name)
		if item.NameEn == "" {
			item.NameEn = strings.TrimSpace(vendorGame.NameEn)
		}
	}
	return item
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
