package game

import (
	"context"
	"strconv"
	"xr-game-server/dao/gamecfgdao"
	"xr-game-server/dto/gamecfgdto"
	"xr-game-server/entity"
	"xr-game-server/module/upload"
)

const (
	appGameListDefaultPageSize = 20
	appGameListMaxPageSize     = 100
)

// GetAppGameCfgList App端分页查询游戏列表(仅已上架,直接读缓存)
func GetAppGameCfgList(_ context.Context, req *gamecfgdto.AppGameCfgListReq) (*gamecfgdto.AppGameCfgListRes, error) {
	page, pageSize := normalizeAppGameListPage(req.Page, req.PageSize)
	all := listOnShelfGameCfgFromCache()
	total := len(all)
	start, end := appGameListPageRange(total, page, pageSize)

	list := make([]*gamecfgdto.AppGameCfgItem, 0, end-start)
	for _, row := range all[start:end] {
		if row == nil {
			continue
		}
		list = append(list, toAppGameCfgItem(row))
	}
	return &gamecfgdto.AppGameCfgListRes{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
	}, nil
}

func listOnShelfGameCfgFromCache() []*entity.GameCfg {
	all := gamecfgdao.GetAllCached()
	list := make([]*entity.GameCfg, 0, len(all))
	for _, row := range all {
		if row == nil || row.Status != entity.GameCfgStatusOnShelf {
			continue
		}
		list = append(list, row)
	}
	return list
}

func toAppGameCfgItem(row *entity.GameCfg) *gamecfgdto.AppGameCfgItem {
	return &gamecfgdto.AppGameCfgItem{
		ID:        strconv.FormatUint(row.ID, 10),
		Name:      row.Name,
		Code:      row.Code,
		LiveCover: upload.GetUrlByName(row.LiveCover),
		Link:      row.Link,
		Sort:      row.Sort,
	}
}

func normalizeAppGameListPage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = appGameListDefaultPageSize
	}
	if pageSize > appGameListMaxPageSize {
		pageSize = appGameListMaxPageSize
	}
	return page, pageSize
}

func appGameListPageRange(total, page, pageSize int) (int, int) {
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	return start, end
}
