package banner

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/bannerdao"
	"xr-game-server/dto/bannerdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
	"xr-game-server/module/upload"
)

func GetBannerList(_ context.Context, req *bannerdto.BannerListReq) (*httpserver.CMSQueryResp, error) {
	total, list := bannerdao.GetBannerList(req)
	for _, row := range list {
		row.ImageName = row.Image
		row.Image = upload.GetUrlByName(row.ImageName)
	}
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}

func GetAppBannerList(_ context.Context, req *bannerdto.AppBannerListReq) (*bannerdto.AppBannerListRes, error) {
	scene := normalizeBannerScene(req.Scene)
	return &bannerdto.AppBannerListRes{List: getAppBannerList(scene)}, nil
}

func normalizeBannerScene(scene uint8) uint8 {
	if scene == entity.HomeBannerSceneLiveRoom {
		return entity.HomeBannerSceneLiveRoom
	}
	return entity.HomeBannerSceneHome
}

func normalizeBannerDirection(scene, direction uint8) uint8 {
	if direction != 0 {
		return direction
	}
	if scene == entity.HomeBannerSceneLiveRoom {
		return entity.HomeBannerPositionLiveHall
	}
	return entity.HomeBannerPositionHomeTop
}

func CreateBanner(_ context.Context, req *bannerdto.CreateBannerReq) (*bannerdto.CreateBannerRes, error) {
	if existing := bannerdao.GetByTitle(req.Title); existing != nil {
		return nil, errercode.CreateCode(errercode.BannerExist)
	}
	scene := normalizeBannerScene(req.Scene)
	row := &entity.HomeBanner{
		Title:     req.Title,
		Image:     req.Image,
		Link:      req.Link,
		Scene:     scene,
		Direction: normalizeBannerDirection(scene, req.Direction),
		Sort:      req.Sort,
		Status:    entity.HomeBannerStatusOffShelf,
	}
	if err := bannerdao.Create(row); err != nil {
		return nil, err
	}
	reloadBannerMemory()
	return &bannerdto.CreateBannerRes{ID: strconv.FormatUint(row.ID, 10)}, nil
}

func UpdateBanner(_ context.Context, req *bannerdto.UpdateBannerReq) (*bannerdto.UpdateBannerRes, error) {
	row := bannerdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.BannerNonExist)
	}
	if existing := bannerdao.GetByTitle(req.Title); existing != nil && existing.ID != req.ID {
		return nil, errercode.CreateCode(errercode.BannerExist)
	}
	row.Title = req.Title
	row.Image = req.Image
	row.Link = req.Link
	row.Scene = normalizeBannerScene(req.Scene)
	row.Direction = normalizeBannerDirection(row.Scene, req.Direction)
	row.Sort = req.Sort
	if err := bannerdao.Update(row); err != nil {
		return nil, err
	}
	reloadBannerMemory()
	return &bannerdto.UpdateBannerRes{Success: true}, nil
}

func DeleteBanner(_ context.Context, req *bannerdto.DeleteBannerReq) (*bannerdto.DeleteBannerRes, error) {
	if row := bannerdao.GetById(req.ID); row == nil {
		return nil, errercode.CreateCode(errercode.BannerNonExist)
	}
	if err := bannerdao.Delete(req.ID); err != nil {
		return nil, err
	}
	reloadBannerMemory()
	return &bannerdto.DeleteBannerRes{Success: true}, nil
}

func OnShelfBanner(_ context.Context, req *bannerdto.OnShelfBannerReq) (*bannerdto.OnShelfBannerRes, error) {
	row := bannerdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.BannerNonExist)
	}
	if row.Status != entity.HomeBannerStatusOnShelf {
		if err := bannerdao.UpdateStatus(req.ID, entity.HomeBannerStatusOnShelf); err != nil {
			return nil, err
		}
		reloadBannerMemory()
	}
	return &bannerdto.OnShelfBannerRes{Success: true, Status: entity.HomeBannerStatusOnShelf}, nil
}

func OffShelfBanner(_ context.Context, req *bannerdto.OffShelfBannerReq) (*bannerdto.OffShelfBannerRes, error) {
	row := bannerdao.GetById(req.ID)
	if row == nil {
		return nil, errercode.CreateCode(errercode.BannerNonExist)
	}
	if row.Status != entity.HomeBannerStatusOffShelf {
		if err := bannerdao.UpdateStatus(req.ID, entity.HomeBannerStatusOffShelf); err != nil {
			return nil, err
		}
		reloadBannerMemory()
	}
	return &bannerdto.OffShelfBannerRes{Success: true, Status: entity.HomeBannerStatusOffShelf}, nil
}
