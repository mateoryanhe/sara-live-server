package game

import (
	"context"

	"xr-game-server/dto/gameplatformdto"
)

const (
	appGameListDefaultPageSize = 20
	appGameListMaxPageSize     = 100
)

// GetAppGameList App 端分页查询已上架游戏列表(读内存).
func GetAppGameList(_ context.Context, req *gameplatformdto.AppGameListReq) (*gameplatformdto.AppGameListRes, error) {
	pageIndex, pageSize := normalizeAppGameListPage(0, 0)
	if req != nil {
		pageIndex, pageSize = normalizeAppGameListPage(req.PageIndex, req.PageSize)
	}
	all := GetAllOnShelfVendorGamesFromMemory()
	total := len(all)
	start, end := appGameListPageRange(total, pageIndex, pageSize)

	list := make([]*gameplatformdto.AppGameListItem, 0, end-start)
	for _, row := range all[start:end] {
		if row == nil {
			continue
		}
		list = append(list, toAppGameListItem(row))
	}
	return &gameplatformdto.AppGameListRes{
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		List:      list,
	}, nil
}

func toAppGameListItem(row *VendorGame) *gameplatformdto.AppGameListItem {
	return &gameplatformdto.AppGameListItem{
		GameCode: row.GameCode,
		NameEn:   row.NameEn,
		Cover:    BuildGameCoverUrl(row.Cover),
		Category: row.Category,
		Platform: row.Platform,
	}
}

func normalizeAppGameListPage(pageIndex, pageSize int) (int, int) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = appGameListDefaultPageSize
	}
	if pageSize > appGameListMaxPageSize {
		pageSize = appGameListMaxPageSize
	}
	return pageIndex, pageSize
}

func appGameListPageRange(total, pageIndex, pageSize int) (int, int) {
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
