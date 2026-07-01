package entryeffect

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/entryeffectdao"
	"xr-game-server/dto/entryeffectdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

func GetEntryEffectList(_ context.Context, req *entryeffectdto.EntryEffectListReq) (*httpserver.CMSQueryResp, error) {
	total, list := entryeffectdao.GetEntryEffectList(req)
	for _, res := range list {
		res.AnimationName = res.Animation
		res.Animation = upload.GetUrlByName(res.AnimationName)
	}
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}

func GetAppEntryEffectList(_ context.Context, _ *entryeffectdto.AppEntryEffectListReq) (*entryeffectdto.AppEntryEffectListRes, error) {
	return &entryeffectdto.AppEntryEffectListRes{List: getAppEntryEffectList()}, nil
}

func CreateEntryEffect(_ context.Context, req *entryeffectdto.CreateEntryEffectReq) (*entryeffectdto.CreateEntryEffectRes, error) {
	if err := validateEntryEffect(req.Name, req.LevelStart, req.LevelEnd, 0); err != nil {
		return nil, err
	}
	if existing := entryeffectdao.GetByName(req.Name); existing != nil {
		return nil, errercode.CreateCode(errercode.EntryEffectExist)
	}

	row := &entity.LiveEntryEffect{
		Name:       req.Name,
		LevelStart: req.LevelStart,
		LevelEnd:   req.LevelEnd,
		Animation:  req.Animation,
		Status:     entity.LiveEntryEffectStatusOffShelf,
	}
	if err := entryeffectdao.Create(row); err != nil {
		return nil, err
	}
	reloadEntryEffectMemory()
	return &entryeffectdto.CreateEntryEffectRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func UpdateEntryEffect(_ context.Context, req *entryeffectdto.UpdateEntryEffectReq) (*entryeffectdto.UpdateEntryEffectRes, error) {
	row := entryeffectdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.EntryEffectNonExist)
	}
	if err := validateEntryEffect(req.Name, req.LevelStart, req.LevelEnd, req.ID); err != nil {
		return nil, err
	}
	if existing := entryeffectdao.GetByName(req.Name); existing != nil && existing.ID != req.ID {
		return nil, errercode.CreateCode(errercode.EntryEffectExist)
	}

	row.Name = req.Name
	row.LevelStart = req.LevelStart
	row.LevelEnd = req.LevelEnd
	row.Animation = req.Animation
	if err := entryeffectdao.Update(row); err != nil {
		return nil, err
	}
	reloadEntryEffectMemory()
	return &entryeffectdto.UpdateEntryEffectRes{Success: true}, nil
}

func DeleteEntryEffect(_ context.Context, req *entryeffectdto.DeleteEntryEffectReq) (*entryeffectdto.DeleteEntryEffectRes, error) {
	if row := entryeffectdao.GetById(req.ID); row == nil {
		return nil, errercode.CreateCode(errercode.EntryEffectNonExist)
	}
	if err := entryeffectdao.Delete(req.ID); err != nil {
		return nil, err
	}
	reloadEntryEffectMemory()
	return &entryeffectdto.DeleteEntryEffectRes{Success: true}, nil
}

func OnShelfEntryEffect(_ context.Context, req *entryeffectdto.OnShelfEntryEffectReq) (*entryeffectdto.OnShelfEntryEffectRes, error) {
	row := entryeffectdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.EntryEffectNonExist)
	}
	if row.Status != entity.LiveEntryEffectStatusOnShelf {
		if err := entryeffectdao.UpdateStatus(req.ID, entity.LiveEntryEffectStatusOnShelf); err != nil {
			return nil, err
		}
		reloadEntryEffectMemory()
	}
	return &entryeffectdto.OnShelfEntryEffectRes{Success: true, Status: entity.LiveEntryEffectStatusOnShelf}, nil
}

func OffShelfEntryEffect(_ context.Context, req *entryeffectdto.OffShelfEntryEffectReq) (*entryeffectdto.OffShelfEntryEffectRes, error) {
	row := entryeffectdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.EntryEffectNonExist)
	}
	if row.Status != entity.LiveEntryEffectStatusOffShelf {
		if err := entryeffectdao.UpdateStatus(req.ID, entity.LiveEntryEffectStatusOffShelf); err != nil {
			return nil, err
		}
		reloadEntryEffectMemory()
	}
	return &entryeffectdto.OffShelfEntryEffectRes{Success: true, Status: entity.LiveEntryEffectStatusOffShelf}, nil
}

func validateEntryEffect(_ string, levelStart, levelEnd int, _ uint64) error {
	if levelEnd < levelStart {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	return nil
}
