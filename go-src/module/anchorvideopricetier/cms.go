package anchorvideopricetier

import (
	"context"
	"strconv"

	"xr-game-server/constants/common"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/anchorvideopricetierdao"
	"xr-game-server/dto/anchorvideopricetierdto"
	entity "xr-game-server/entity/live"
	"xr-game-server/errercode"
)

func GetPriceTierList(_ context.Context, req *anchorvideopricetierdto.PriceTierListReq) (*httpserver.CMSQueryResp, error) {
	total, list := anchorvideopricetierdao.GetList(req)
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}

func CreatePriceTier(_ context.Context, req *anchorvideopricetierdto.CreatePriceTierReq) (*anchorvideopricetierdto.CreatePriceTierRes, error) {
	if req == nil || req.Price <= common.Zero {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := &entity.AnchorVideoPriceTier{
		Price:  req.Price,
		Status: entity.AnchorVideoPriceTierStatusOnShelf,
	}
	if err := anchorvideopricetierdao.Create(row); err != nil {
		return nil, err
	}
	reloadMemory()
	return &anchorvideopricetierdto.CreatePriceTierRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func UpdatePriceTier(_ context.Context, req *anchorvideopricetierdto.UpdatePriceTierReq) (*anchorvideopricetierdto.UpdatePriceTierRes, error) {
	if req == nil || req.Price <= common.Zero {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := anchorvideopricetierdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.AnchorVideoPriceTierNonExist)
	}
	row.Price = req.Price
	if err := anchorvideopricetierdao.Update(row); err != nil {
		return nil, err
	}
	reloadMemory()
	return &anchorvideopricetierdto.UpdatePriceTierRes{Success: true}, nil
}

func DeletePriceTier(_ context.Context, req *anchorvideopricetierdto.DeletePriceTierReq) (*anchorvideopricetierdto.DeletePriceTierRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if row := anchorvideopricetierdao.GetById(req.ID); row == nil {
		return nil, errercode.CreateCode(errercode.AnchorVideoPriceTierNonExist)
	}
	if err := anchorvideopricetierdao.Delete(req.ID); err != nil {
		return nil, err
	}
	reloadMemory()
	return &anchorvideopricetierdto.DeletePriceTierRes{Success: true}, nil
}

func OnShelfPriceTier(_ context.Context, req *anchorvideopricetierdto.OnShelfPriceTierReq) (*anchorvideopricetierdto.OnShelfPriceTierRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := anchorvideopricetierdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.AnchorVideoPriceTierNonExist)
	}
	if row.Status != entity.AnchorVideoPriceTierStatusOnShelf {
		if err := anchorvideopricetierdao.UpdateStatus(req.ID, entity.AnchorVideoPriceTierStatusOnShelf); err != nil {
			return nil, err
		}
		reloadMemory()
	}
	return &anchorvideopricetierdto.OnShelfPriceTierRes{Success: true, Status: entity.AnchorVideoPriceTierStatusOnShelf}, nil
}

func OffShelfPriceTier(_ context.Context, req *anchorvideopricetierdto.OffShelfPriceTierReq) (*anchorvideopricetierdto.OffShelfPriceTierRes, error) {
	if req == nil {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	row := anchorvideopricetierdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.AnchorVideoPriceTierNonExist)
	}
	if row.Status != entity.AnchorVideoPriceTierStatusOffShelf {
		if err := anchorvideopricetierdao.UpdateStatus(req.ID, entity.AnchorVideoPriceTierStatusOffShelf); err != nil {
			return nil, err
		}
		reloadMemory()
	}
	return &anchorvideopricetierdto.OffShelfPriceTierRes{Success: true, Status: entity.AnchorVideoPriceTierStatusOffShelf}, nil
}
