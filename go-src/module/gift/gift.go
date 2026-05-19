package gift

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/giftdao"
	"xr-game-server/dto/giftdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

// GetGiftList 分页获取礼物列表(CMS)
func GetGiftList(_ context.Context, req *giftdto.GiftListReq) (*httpserver.CMSQueryResp, error) {
	total, list := giftdao.GetGiftList(req)
	return &httpserver.CMSQueryResp{
		Total: total,
		Data:  list,
	}, nil
}

// GetAppGiftList App端查询礼物列表(仅返回已上架礼物,带缓存)
// 缓存以礼物 ID 为主键,在礼物创建/修改/删除/上下架时清空,下次访问时重新加载。
func GetAppGiftList(_ context.Context, _ *giftdto.AppGiftListReq) (*giftdto.AppGiftListRes, error) {
	return &giftdto.AppGiftListRes{List: getGiftCache()}, nil
}

// CreateGift 创建礼物
func CreateGift(_ context.Context, req *giftdto.CreateGiftReq) (*giftdto.CreateGiftRes, error) {
	if existing := giftdao.GetGiftByName(req.Name); existing != nil {
		return nil, errercode.CreateCode(errercode.GiftExist)
	}

	g := &entity.LiveGift{
		Name:        req.Name,
		Icon:        req.Icon,
		Animation:   req.Animation,
		Price:       req.Price,
		Category:    req.Category,
		Sort:        req.Sort,
		Status:      entity.LiveGiftStatusOffShelf, // 新建默认下架,需手动上架
		Description: req.Description,
	}
	if err := giftdao.CreateGift(g); err != nil {
		return nil, err
	}
	invalidateGiftCache()
	return &giftdto.CreateGiftRes{ID: strconv.FormatUint(g.ID, 10)}, nil
}

// UpdateGift 修改礼物(不修改上下架状态)
func UpdateGift(_ context.Context, req *giftdto.UpdateGiftReq) (*giftdto.UpdateGiftRes, error) {
	g := giftdao.GetGiftById(req.ID)
	if g == nil {
		return nil, errercode.CreateCode(errercode.GiftNonExist)
	}
	if existing := giftdao.GetGiftByName(req.Name); existing != nil && existing.ID != req.ID {
		return nil, errercode.CreateCode(errercode.GiftExist)
	}

	g.Name = req.Name
	g.Icon = req.Icon
	g.Animation = req.Animation
	g.Price = req.Price
	g.Category = req.Category
	g.Sort = req.Sort
	g.Description = req.Description

	if err := giftdao.UpdateGift(g); err != nil {
		return nil, err
	}
	invalidateGiftCache()
	return &giftdto.UpdateGiftRes{Success: true}, nil
}

// DeleteGift 删除礼物
func DeleteGift(_ context.Context, req *giftdto.DeleteGiftReq) (*giftdto.DeleteGiftRes, error) {
	if g := giftdao.GetGiftById(req.ID); g == nil {
		return nil, errercode.CreateCode(errercode.GiftNonExist)
	}
	if err := giftdao.DeleteGift(req.ID); err != nil {
		return nil, err
	}
	invalidateGiftCache()
	return &giftdto.DeleteGiftRes{Success: true}, nil
}

// OnShelfGift 上架礼物
func OnShelfGift(_ context.Context, req *giftdto.OnShelfGiftReq) (*giftdto.OnShelfGiftRes, error) {
	g := giftdao.GetGiftById(req.ID)
	if g == nil {
		return nil, errercode.CreateCode(errercode.GiftNonExist)
	}
	if g.Status != entity.LiveGiftStatusOnShelf {
		if err := giftdao.UpdateGiftStatus(req.ID, entity.LiveGiftStatusOnShelf); err != nil {
			return nil, err
		}
		invalidateGiftCache()
	}
	return &giftdto.OnShelfGiftRes{Success: true, Status: entity.LiveGiftStatusOnShelf}, nil
}

// OffShelfGift 下架礼物
func OffShelfGift(_ context.Context, req *giftdto.OffShelfGiftReq) (*giftdto.OffShelfGiftRes, error) {
	g := giftdao.GetGiftById(req.ID)
	if g == nil {
		return nil, errercode.CreateCode(errercode.GiftNonExist)
	}
	if g.Status != entity.LiveGiftStatusOffShelf {
		if err := giftdao.UpdateGiftStatus(req.ID, entity.LiveGiftStatusOffShelf); err != nil {
			return nil, err
		}
		invalidateGiftCache()
	}
	return &giftdto.OffShelfGiftRes{Success: true, Status: entity.LiveGiftStatusOffShelf}, nil
}
