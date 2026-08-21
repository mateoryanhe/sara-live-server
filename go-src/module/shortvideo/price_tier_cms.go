package shortvideo

import (
	"context"
	"strconv"
	"xr-game-server/constants/common"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity/shortvideo"
	"xr-game-server/errercode"
)

func GetShortVideoPriceTierList(_ context.Context, req *shortvideodto.ShortVideoPriceTierListReq) (*httpserver.CMSQueryResp, error) {
	total, list := shortvideodao.GetPriceTierList(req)
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}

func GetAppShortVideoPriceTierList(_ context.Context, _ *shortvideodto.AppShortVideoPriceTierListReq) (*shortvideodto.AppShortVideoPriceTierListRes, error) {
	return &shortvideodto.AppShortVideoPriceTierListRes{List: getAppPriceTierList()}, nil
}

func CreateShortVideoPriceTier(_ context.Context, req *shortvideodto.CreateShortVideoPriceTierReq) (*shortvideodto.CreateShortVideoPriceTierRes, error) {
	if req.Price <= common.Zero {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := &entity.ShortVideoPriceTier{
		Price:  req.Price,
		Status: entity.ShortVideoPriceTierStatusOnShelf,
	}
	if err := shortvideodao.CreatePriceTier(row); err != nil {
		return nil, err
	}
	reloadPriceTierMemory()
	return &shortvideodto.CreateShortVideoPriceTierRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func UpdateShortVideoPriceTier(_ context.Context, req *shortvideodto.UpdateShortVideoPriceTierReq) (*shortvideodto.UpdateShortVideoPriceTierRes, error) {
	if req.Price <= common.Zero {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := shortvideodao.GetPriceTierById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoPriceTierNonExist)
	}
	row.Price = req.Price
	if err := shortvideodao.UpdatePriceTier(row); err != nil {
		return nil, err
	}
	reloadPriceTierMemory()
	return &shortvideodto.UpdateShortVideoPriceTierRes{Success: true}, nil
}

func DeleteShortVideoPriceTier(_ context.Context, req *shortvideodto.DeleteShortVideoPriceTierReq) (*shortvideodto.DeleteShortVideoPriceTierRes, error) {
	if row := shortvideodao.GetPriceTierById(req.ID); row == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoPriceTierNonExist)
	}
	if err := shortvideodao.DeletePriceTier(req.ID); err != nil {
		return nil, err
	}
	reloadPriceTierMemory()
	return &shortvideodto.DeleteShortVideoPriceTierRes{Success: true}, nil
}

func OnShelfShortVideoPriceTier(_ context.Context, req *shortvideodto.OnShelfShortVideoPriceTierReq) (*shortvideodto.OnShelfShortVideoPriceTierRes, error) {
	row := shortvideodao.GetPriceTierById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoPriceTierNonExist)
	}
	if row.Status != entity.ShortVideoPriceTierStatusOnShelf {
		if err := shortvideodao.UpdatePriceTierStatus(req.ID, entity.ShortVideoPriceTierStatusOnShelf); err != nil {
			return nil, err
		}
		reloadPriceTierMemory()
	}
	return &shortvideodto.OnShelfShortVideoPriceTierRes{Success: true, Status: entity.ShortVideoPriceTierStatusOnShelf}, nil
}

func OffShelfShortVideoPriceTier(_ context.Context, req *shortvideodto.OffShelfShortVideoPriceTierReq) (*shortvideodto.OffShelfShortVideoPriceTierRes, error) {
	row := shortvideodao.GetPriceTierById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoPriceTierNonExist)
	}
	if row.Status != entity.ShortVideoPriceTierStatusOffShelf {
		if err := shortvideodao.UpdatePriceTierStatus(req.ID, entity.ShortVideoPriceTierStatusOffShelf); err != nil {
			return nil, err
		}
		reloadPriceTierMemory()
	}
	return &shortvideodto.OffShelfShortVideoPriceTierRes{Success: true, Status: entity.ShortVideoPriceTierStatusOffShelf}, nil
}
