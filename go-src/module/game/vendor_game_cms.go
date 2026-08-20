package game

import (
	"context"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/gameplatformdto"
)

const (
	cmsVendorGameDefaultPageSize = 20
	cmsVendorGameMaxPageSize     = 100
)

// GetVendorGameList CMS 分页查询游戏库(读 vendor_game_libs 表).
func GetVendorGameList(_ context.Context, req *gameplatformdto.VendorGameListReq) (*httpserver.CMSQueryResp, error) {
	pageIndex, pageSize := normalizeCMSVendorGamePage(1, cmsVendorGameDefaultPageSize)
	q := &cfgdao.VendorGameLibQuery{
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}
	if req != nil {
		pageIndex, pageSize = normalizeCMSVendorGamePage(req.PageIndex, req.PageSize)
		q.PageIndex = pageIndex
		q.PageSize = pageSize
		q.GameCode = req.GameCode
		q.Name = req.Name
		q.Platform = req.Platform
		q.Category = req.Category
	}

	total, rows := cfgdao.QueryVendorGameLibs(q)
	list := make([]*gameplatformdto.VendorGameListItem, 0, len(rows))
	shelfSet := cfgdao.GetGameCfgShelfKeySetFromMemory()
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, toVendorGameListItem(toVendorGame(row), shelfSet))
	}
	return httpserver.NewCMSQueryResp(total, list), nil
}

// ReloadVendorGameCacheCMS 从第三方全量同步游戏库表.
func ReloadVendorGameCacheCMS(ctx context.Context, _ *gameplatformdto.ReloadVendorGameCacheReq) (*gameplatformdto.ReloadVendorGameCacheRes, error) {
	count, err := SyncVendorGameLibraryFromVendor(ctx)
	if err != nil {
		return nil, err
	}
	return &gameplatformdto.ReloadVendorGameCacheRes{
		Success: true,
		Count:   count,
	}, nil
}

func toVendorGameListItem(row *VendorGame, shelfSet map[string]struct{}) *gameplatformdto.VendorGameListItem {
	if row == nil {
		return &gameplatformdto.VendorGameListItem{}
	}
	_, onShelf := shelfSet[cfgdao.GameCfgShelfKey(row.GameCode, row.Platform)]
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
