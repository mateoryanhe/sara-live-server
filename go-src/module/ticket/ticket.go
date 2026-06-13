package ticket

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/ticketdao"
	"xr-game-server/dto/ticketdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func GetTicketList(_ context.Context, req *ticketdto.TicketListReq) (*httpserver.CMSQueryResp, error) {
	total, list := ticketdao.GetTicketList(req)
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}

func GetAppTicketList(_ context.Context, _ *ticketdto.AppTicketListReq) (*ticketdto.AppTicketListRes, error) {
	return &ticketdto.AppTicketListRes{List: getAppTicketList()}, nil
}

func CreateTicket(_ context.Context, req *ticketdto.CreateTicketReq) (*ticketdto.CreateTicketRes, error) {
	row := &entity.LiveTicket{
		Price:  req.Price,
		Sort:   req.Sort,
		Status: entity.LiveTicketStatusOnShelf,
	}
	if err := ticketdao.Create(row); err != nil {
		return nil, err
	}
	reloadTicketMemory()
	return &ticketdto.CreateTicketRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func UpdateTicket(_ context.Context, req *ticketdto.UpdateTicketReq) (*ticketdto.UpdateTicketRes, error) {
	row := ticketdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.TicketNonExist)
	}
	row.Price = req.Price
	row.Sort = req.Sort
	if err := ticketdao.Update(row); err != nil {
		return nil, err
	}
	reloadTicketMemory()
	return &ticketdto.UpdateTicketRes{Success: true}, nil
}

func DeleteTicket(_ context.Context, req *ticketdto.DeleteTicketReq) (*ticketdto.DeleteTicketRes, error) {
	if row := ticketdao.GetById(req.ID); row == nil {
		return nil, errercode.CreateCode(errercode.TicketNonExist)
	}
	if err := ticketdao.Delete(req.ID); err != nil {
		return nil, err
	}
	reloadTicketMemory()
	return &ticketdto.DeleteTicketRes{Success: true}, nil
}

func OnShelfTicket(_ context.Context, req *ticketdto.OnShelfTicketReq) (*ticketdto.OnShelfTicketRes, error) {
	row := ticketdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.TicketNonExist)
	}
	if row.Status != entity.LiveTicketStatusOnShelf {
		if err := ticketdao.UpdateStatus(req.ID, entity.LiveTicketStatusOnShelf); err != nil {
			return nil, err
		}
		reloadTicketMemory()
	}
	return &ticketdto.OnShelfTicketRes{Success: true, Status: entity.LiveTicketStatusOnShelf}, nil
}

func OffShelfTicket(_ context.Context, req *ticketdto.OffShelfTicketReq) (*ticketdto.OffShelfTicketRes, error) {
	row := ticketdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.TicketNonExist)
	}
	if row.Status != entity.LiveTicketStatusOffShelf {
		if err := ticketdao.UpdateStatus(req.ID, entity.LiveTicketStatusOffShelf); err != nil {
			return nil, err
		}
		reloadTicketMemory()
	}
	return &ticketdto.OffShelfTicketRes{Success: true, Status: entity.LiveTicketStatusOffShelf}, nil
}
