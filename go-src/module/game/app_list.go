package game

import (
	"context"

	"xr-game-server/dto/gameplatformdto"
	"xr-game-server/entity/game"
)

const (
	appGameListDefaultPageSize = 20
	appGameListMaxPageSize     = 100
)

// GetAppGameList App 端分页查询已上架游戏(读 game_cfgs 永久缓存).
func GetAppGameList(_ context.Context, req *gameplatformdto.AppGameListReq) (*gameplatformdto.AppGameListRes, error) {
	pageIndex, pageSize := normalizeAppGameListPage(0, 0)
	if req != nil {
		pageIndex, pageSize = normalizeAppGameListPage(req.PageIndex, req.PageSize)
	}
	all := GetAllOnShelfGamesFromMemory()
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

func toAppGameListItem(row *entity.GameCfg) *gameplatformdto.AppGameListItem {
	return &gameplatformdto.AppGameListItem{
		GameCode: row.GameCode,
		NameEn:   ResolveAppGameName(row),
		Cover:    ResolveAppGameCover(row),
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
