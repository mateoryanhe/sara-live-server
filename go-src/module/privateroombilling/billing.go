package privateroombilling

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/privateroombillingdao"
	"xr-game-server/dto/privateroombillingdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func GetBillingList(_ context.Context, req *privateroombillingdto.BillingListReq) (*httpserver.CMSQueryResp, error) {
	total, list := privateroombillingdao.GetBillingList(req)
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}

func GetAppBillingList(_ context.Context, _ *privateroombillingdto.AppBillingListReq) (*privateroombillingdto.AppBillingListRes, error) {
	return &privateroombillingdto.AppBillingListRes{List: getAppBillingList()}, nil
}

func CreateBilling(_ context.Context, req *privateroombillingdto.CreateBillingReq) (*privateroombillingdto.CreateBillingRes, error) {
	row := &entity.LivePrivateRoomBilling{
		PricePerMinute: req.PricePerMinute,
		Sort:           req.Sort,
		Status:         entity.LivePrivateRoomBillingStatusOnShelf,
	}
	if err := privateroombillingdao.Create(row); err != nil {
		return nil, err
	}
	reloadBillingMemory()
	return &privateroombillingdto.CreateBillingRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func UpdateBilling(_ context.Context, req *privateroombillingdto.UpdateBillingReq) (*privateroombillingdto.UpdateBillingRes, error) {
	row := privateroombillingdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.PrivateRoomBillingNonExist)
	}
	row.PricePerMinute = req.PricePerMinute
	row.Sort = req.Sort
	if err := privateroombillingdao.Update(row); err != nil {
		return nil, err
	}
	reloadBillingMemory()
	return &privateroombillingdto.UpdateBillingRes{Success: true}, nil
}

func DeleteBilling(_ context.Context, req *privateroombillingdto.DeleteBillingReq) (*privateroombillingdto.DeleteBillingRes, error) {
	if row := privateroombillingdao.GetById(req.ID); row == nil {
		return nil, errercode.CreateCode(errercode.PrivateRoomBillingNonExist)
	}
	if err := privateroombillingdao.Delete(req.ID); err != nil {
		return nil, err
	}
	reloadBillingMemory()
	return &privateroombillingdto.DeleteBillingRes{Success: true}, nil
}

func OnShelfBilling(_ context.Context, req *privateroombillingdto.OnShelfBillingReq) (*privateroombillingdto.OnShelfBillingRes, error) {
	row := privateroombillingdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.PrivateRoomBillingNonExist)
	}
	if row.Status != entity.LivePrivateRoomBillingStatusOnShelf {
		if err := privateroombillingdao.UpdateStatus(req.ID, entity.LivePrivateRoomBillingStatusOnShelf); err != nil {
			return nil, err
		}
		reloadBillingMemory()
	}
	return &privateroombillingdto.OnShelfBillingRes{Success: true, Status: entity.LivePrivateRoomBillingStatusOnShelf}, nil
}

func OffShelfBilling(_ context.Context, req *privateroombillingdto.OffShelfBillingReq) (*privateroombillingdto.OffShelfBillingRes, error) {
	row := privateroombillingdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.PrivateRoomBillingNonExist)
	}
	if row.Status != entity.LivePrivateRoomBillingStatusOffShelf {
		if err := privateroombillingdao.UpdateStatus(req.ID, entity.LivePrivateRoomBillingStatusOffShelf); err != nil {
			return nil, err
		}
		reloadBillingMemory()
	}
	return &privateroombillingdto.OffShelfBillingRes{Success: true, Status: entity.LivePrivateRoomBillingStatusOffShelf}, nil
}
