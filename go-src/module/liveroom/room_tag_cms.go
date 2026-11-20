package liveroom

import (
	"context"
	"strconv"
	"strings"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/liveroomdao"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity/live"
	"xr-game-server/errercode"
)

func GetLiveRoomTagList(_ context.Context, req *liveroomdto.LiveRoomTagListReq) (*httpserver.CMSQueryResp, error) {
	total, list := liveroomdao.GetRoomTagList(req)
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}

func GetAppLiveRoomTagList(_ context.Context, _ *liveroomdto.AppLiveRoomTagListReq) (*liveroomdto.AppLiveRoomTagListRes, error) {
	return &liveroomdto.AppLiveRoomTagListRes{List: getAppRoomTagList()}, nil
}

func GetAppLiveRoomNormalTagList(_ context.Context, _ *liveroomdto.AppLiveRoomNormalTagListReq) (*liveroomdto.AppLiveRoomTagListRes, error) {
	return &liveroomdto.AppLiveRoomTagListRes{List: getAppRoomTagNormalList()}, nil
}

func normalizeLiveRoomTagCode(code string) string {
	return strings.TrimSpace(code)
}

func validateLiveRoomTagCodeUnique(code string, excludeID uint64) error {
	code = normalizeLiveRoomTagCode(code)
	if code == "" {
		return nil
	}
	if existing := liveroomdao.GetRoomTagByCode(code); existing != nil && existing.ID != excludeID {
		return errercode.CreateCode(errercode.LiveRoomTagExist)
	}
	return nil
}

func CreateLiveRoomTag(_ context.Context, req *liveroomdto.CreateLiveRoomTagReq) (*liveroomdto.CreateLiveRoomTagRes, error) {
	if existing := liveroomdao.GetRoomTagByName(req.Name); existing != nil {
		return nil, errercode.CreateCode(errercode.LiveRoomTagExist)
	}
	if err := validateLiveRoomTagCodeUnique(req.Code, 0); err != nil {
		return nil, err
	}
	row := &entity.LiveRoomTag{
		Name:      req.Name,
		Code:      normalizeLiveRoomTagCode(req.Code),
		Sort:      req.Sort,
		IsSpecial: req.IsSpecial,
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
	if err := validateLiveRoomTagCodeUnique(req.Code, req.ID); err != nil {
		return nil, err
	}
	row.Name = req.Name
	row.Code = normalizeLiveRoomTagCode(req.Code)
	row.Sort = req.Sort
	row.IsSpecial = req.IsSpecial
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
