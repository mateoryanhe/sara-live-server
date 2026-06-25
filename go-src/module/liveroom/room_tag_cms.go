package liveroom

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func GetLiveRoomTagList(_ context.Context, req *liveroomdto.LiveRoomTagListReq) (*httpserver.CMSQueryResp, error) {
	total, list := liveroomdao.GetRoomTagList(req)
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}

func GetAppLiveRoomTagList(_ context.Context, _ *liveroomdto.AppLiveRoomTagListReq) (*liveroomdto.AppLiveRoomTagListRes, error) {
	return &liveroomdto.AppLiveRoomTagListRes{List: getAppRoomTagList()}, nil
}

func CreateLiveRoomTag(_ context.Context, req *liveroomdto.CreateLiveRoomTagReq) (*liveroomdto.CreateLiveRoomTagRes, error) {
	if existing := liveroomdao.GetRoomTagByName(req.Name); existing != nil {
		return nil, errercode.CreateCode(errercode.LiveRoomTagExist)
	}
	row := &entity.LiveRoomTag{
		Name: req.Name,
		Sort: req.Sort,
	}
	if err := liveroomdao.CreateRoomTag(row); err != nil {
		return nil, err
	}
	reloadRoomTagMemory()
	return &liveroomdto.CreateLiveRoomTagRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func UpdateLiveRoomTag(_ context.Context, req *liveroomdto.UpdateLiveRoomTagReq) (*liveroomdto.UpdateLiveRoomTagRes, error) {
	row := liveroomdao.GetRoomTagById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomTagNonExist)
	}
	if existing := liveroomdao.GetRoomTagByName(req.Name); existing != nil && existing.ID != req.ID {
		return nil, errercode.CreateCode(errercode.LiveRoomTagExist)
	}
	row.Name = req.Name
	row.Sort = req.Sort
	if err := liveroomdao.UpdateRoomTag(row); err != nil {
		return nil, err
	}
	reloadRoomTagMemory()
	return &liveroomdto.UpdateLiveRoomTagRes{Success: true}, nil
}

func DeleteLiveRoomTag(_ context.Context, req *liveroomdto.DeleteLiveRoomTagReq) (*liveroomdto.DeleteLiveRoomTagRes, error) {
	if row := liveroomdao.GetRoomTagById(req.ID); row == nil {
		return nil, errercode.CreateCode(errercode.LiveRoomTagNonExist)
	}
	if err := liveroomdao.DeleteRoomTag(req.ID); err != nil {
		return nil, err
	}
	reloadRoomTagMemory()
	return &liveroomdto.DeleteLiveRoomTagRes{Success: true}, nil
}

func validateLiveRoomTag(tagId uint64) error {
	if tagId == 0 {
		return errercode.CreateCode(errercode.InvalidParam)
	}
	if GetRoomTagFromMemoryById(tagId) == nil {
		return errercode.CreateCode(errercode.LiveRoomTagNonExist)
	}
	return nil
}
