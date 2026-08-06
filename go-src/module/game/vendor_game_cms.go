package game

import (
	"context"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/gameplatformdto"
)

const (
	cmsVendorGameDefaultPageSize = 20
	cmsVendorGameMaxPageSize     = 100
)

// GetVendorGameList CMS 分页查询第三方游戏(搜索时可全量拉取并覆盖 30 分钟浏览缓存).
func GetVendorGameList(ctx context.Context, req *gameplatformdto.VendorGameListReq) (*httpserver.CMSQueryResp, error) {
	if req != nil && req.RefreshFromVendor {
		if err := ForceRefreshVendorBrowseCache(ctx); err != nil {
			return nil, err
		}
	}
	all := filterVendorGames(GetAllVendorBrowseGamesFromMemory(), req)
	total := len(all)
	pageIndex, pageSize := normalizeCMSVendorGamePage(req.PageIndex, req.PageSize)
	start, end := cmsVendorGamePageRange(total, pageIndex, pageSize)

	list := make([]*gameplatformdto.VendorGameListItem, 0, end-start)
	shelfSet := cfgdao.GetGameCfgCodeSetFromMemory()
	for _, row := range all[start:end] {
		if row == nil {
			continue
		}
		list = append(list, toVendorGameListItem(row, shelfSet))
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

// ReloadVendorGameCacheCMS 从第三方重新拉取游戏列表到浏览缓存.
func ReloadVendorGameCacheCMS(ctx context.Context, _ *gameplatformdto.ReloadVendorGameCacheReq) (*gameplatformdto.ReloadVendorGameCacheRes, error) {
	if err := ForceRefreshVendorBrowseCache(ctx); err != nil {
		return nil, err
	}
	return &gameplatformdto.ReloadVendorGameCacheRes{
		Success: true,
		Count:   len(GetAllVendorBrowseGamesFromMemory()),
	}, nil
}

func filterVendorGames(all []*VendorGame, req *gameplatformdto.VendorGameListReq) []*VendorGame {
	if req == nil {
		return all
	}
	gameCode := strings.TrimSpace(req.GameCode)
	name := strings.TrimSpace(req.Name)
	platform := strings.TrimSpace(req.Platform)
	category := strings.TrimSpace(req.Category)
	if gameCode == "" && name == "" && platform == "" && category == "" {
		return all
	}

	list := make([]*VendorGame, 0, len(all))
	for _, row := range all {
		if row == nil {
			continue
		}
		if gameCode != "" && !strings.Contains(strings.ToLower(row.GameCode), strings.ToLower(gameCode)) {
			continue
		}
		if name != "" && !strings.Contains(strings.ToLower(row.Name), strings.ToLower(name)) &&
			!strings.Contains(strings.ToLower(row.NameEn), strings.ToLower(name)) {
			continue
		}
		if platform != "" && !strings.Contains(strings.ToLower(row.Platform), strings.ToLower(platform)) {
			continue
		}
		if category != "" && !strings.Contains(strings.ToLower(row.Category), strings.ToLower(category)) {
			continue
		}
		list = append(list, row)
	}
	return list
}

func toVendorGameListItem(row *VendorGame, shelfSet map[string]struct{}) *gameplatformdto.VendorGameListItem {
	_, onShelf := shelfSet[row.GameCode]
	return &gameplatformdto.VendorGameListItem{
		GameCode: row.GameCode,
		Name:     row.Name,
		NameEn:   row.NameEn,
		Category: row.Category,
		Cover:    BuildGameCoverUrl(row.Cover),
		Platform: row.Platform,
		OnShelf:  onShelf,
	}
}

func normalizeCMSVendorGamePage(pageIndex, pageSize int) (int, int) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = cmsVendorGameDefaultPageSize
	}
	if pageSize > cmsVendorGameMaxPageSize {
		pageSize = cmsVendorGameMaxPageSize
	}
	return pageIndex, pageSize
}

func cmsVendorGamePageRange(total, pageIndex, pageSize int) (int, int) {
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
