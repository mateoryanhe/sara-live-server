package shortvideo

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/shortvideodao"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

func GetShortVideoCategoryList(_ context.Context, req *shortvideodto.ShortVideoCategoryListReq) (*httpserver.CMSQueryResp, error) {
	total, list := shortvideodao.GetCategoryList(req)
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}

func GetAppShortVideoCategoryList(_ context.Context, _ *shortvideodto.AppShortVideoCategoryListReq) (*shortvideodto.AppShortVideoCategoryListRes, error) {
	return &shortvideodto.AppShortVideoCategoryListRes{List: getAppCategoryList()}, nil
}

func CreateShortVideoCategory(_ context.Context, req *shortvideodto.CreateShortVideoCategoryReq) (*shortvideodto.CreateShortVideoCategoryRes, error) {
	if existing := shortvideodao.GetCategoryByName(req.Name); existing != nil {
		return nil, errercode.CreateCode(errercode.ShortVideoCategoryExist)
	}
	row := &entity.ShortVideoCategory{
		Name: req.Name,
		Sort: req.Sort,
	}
	if err := shortvideodao.CreateCategory(row); err != nil {
		return nil, err
	}
	reloadCategoryMemory()
	return &shortvideodto.CreateShortVideoCategoryRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func UpdateShortVideoCategory(_ context.Context, req *shortvideodto.UpdateShortVideoCategoryReq) (*shortvideodto.UpdateShortVideoCategoryRes, error) {
	row := shortvideodao.GetCategoryById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoCategoryNonExist)
	}
	if existing := shortvideodao.GetCategoryByName(req.Name); existing != nil && existing.ID != req.ID {
		return nil, errercode.CreateCode(errercode.ShortVideoCategoryExist)
	}
	row.Name = req.Name
	row.Sort = req.Sort
	if err := shortvideodao.UpdateCategory(row); err != nil {
		return nil, err
	}
	reloadCategoryMemory()
	return &shortvideodto.UpdateShortVideoCategoryRes{Success: true}, nil
}

func DeleteShortVideoCategory(_ context.Context, req *shortvideodto.DeleteShortVideoCategoryReq) (*shortvideodto.DeleteShortVideoCategoryRes, error) {
	if row := shortvideodao.GetCategoryById(req.ID); row == nil {
		return nil, errercode.CreateCode(errercode.ShortVideoCategoryNonExist)
	}
	if err := shortvideodao.DeleteCategory(req.ID); err != nil {
		return nil, err
	}
	reloadCategoryMemory()
	return &shortvideodto.DeleteShortVideoCategoryRes{Success: true}, nil
}
