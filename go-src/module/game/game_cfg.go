package game

import (
	"context"
	"strconv"
	"strings"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/gamecfgdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

func GetList(_ context.Context, req *gamecfgdto.GameCfgListReq) (*httpserver.CMSQueryResp, error) {
	total, list := cfgdao.GetGameCfgList(req)
	for _, item := range list {
		if item == nil {
			continue
		}
		item.LiveCoverUrl = upload.GetUrlByName(item.LiveCover)
	}
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}

func Create(_ context.Context, req *gamecfgdto.CreateGameCfgReq) (*gamecfgdto.CreateGameCfgRes, error) {
	code := strings.TrimSpace(req.Code)
	if existing := cfgdao.GetGameCfgByCode(code); existing != nil {
		return nil, errercode.CreateCode(errercode.GameCfgExist)
	}
	row := &entity.GameCfg{
		Name:      strings.TrimSpace(req.Name),
		Code:      code,
		LiveCover: strings.TrimSpace(req.LiveCover),
		Link:      strings.TrimSpace(req.Link),
		Sort:      req.Sort,
	}
	if err := cfgdao.CreateGameCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadGameCfgCache()
	return &gamecfgdto.CreateGameCfgRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func Update(_ context.Context, req *gamecfgdto.UpdateGameCfgReq) (*gamecfgdto.UpdateGameCfgRes, error) {
	row := cfgdao.GetGameCfgById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.GameCfgNonExist)
	}
	code := strings.TrimSpace(req.Code)
	if existing := cfgdao.GetGameCfgByCode(code); existing != nil && existing.ID != req.ID {
		return nil, errercode.CreateCode(errercode.GameCfgExist)
	}
	row.Name = strings.TrimSpace(req.Name)
	row.Code = code
	row.LiveCover = strings.TrimSpace(req.LiveCover)
	row.Link = strings.TrimSpace(req.Link)
	row.Sort = req.Sort
	if err := cfgdao.UpdateGameCfg(row); err != nil {
		return nil, err
	}
	cfgdao.ReloadGameCfgCache()
	return &gamecfgdto.UpdateGameCfgRes{Success: true}, nil
}

func Delete(_ context.Context, req *gamecfgdto.DeleteGameCfgReq) (*gamecfgdto.DeleteGameCfgRes, error) {
	if row := cfgdao.GetGameCfgById(req.ID); row == nil {
		return nil, errercode.CreateCode(errercode.GameCfgNonExist)
	}
	if err := cfgdao.DeleteGameCfg(req.ID); err != nil {
		return nil, err
	}
	cfgdao.ReloadGameCfgCache()
	return &gamecfgdto.DeleteGameCfgRes{Success: true}, nil
}
