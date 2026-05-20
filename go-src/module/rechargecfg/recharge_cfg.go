package rechargecfg

import (
	"context"
	"strconv"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/rechargecfgdao"
	"xr-game-server/dto/rechargecfgdto"
	"xr-game-server/entity"
	"xr-game-server/errercode"
)

// ===== CMS =====

// GetList CMS分页查询(全部状态)
func GetList(_ context.Context, req *rechargecfgdto.RechargeCfgListReq) (*httpserver.CMSQueryResp, error) {
	total, list := rechargecfgdao.GetList(req)
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}

// Create 创建充值配置(默认下架,需手动上架)
func Create(_ context.Context, req *rechargecfgdto.CreateRechargeCfgReq) (*rechargecfgdto.CreateRechargeCfgRes, error) {
	if existing := rechargecfgdao.GetByName(req.Name); existing != nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgExist)
	}
	currency := req.Currency
	if currency == "" {
		currency = "CNY"
	}
	cfg := &entity.RechargeCfg{
		Name:         req.Name,
		Icon:         req.Icon,
		Diamond:      req.Diamond,
		ExtraDiamond: req.ExtraDiamond,
		Price:        req.Price,
		Currency:     currency,
		ProductId:    req.ProductId,
		Sort:         req.Sort,
		Status:       entity.RechargeCfgStatusOffShelf,
		Description:  req.Description,
	}
	if err := rechargecfgdao.Create(cfg); err != nil {
		return nil, err
	}
	reloadRechargeCfgCache()
	return &rechargecfgdto.CreateRechargeCfgRes{ID: strconv.FormatUint(cfg.ID, 10)}, nil
}

// Update 修改充值配置(不修改上下架状态)
func Update(_ context.Context, req *rechargecfgdto.UpdateRechargeCfgReq) (*rechargecfgdto.UpdateRechargeCfgRes, error) {
	cfg := rechargecfgdao.GetById(req.ID)
	if cfg == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	if existing := rechargecfgdao.GetByName(req.Name); existing != nil && existing.ID != req.ID {
		return nil, errercode.CreateCode(errercode.RechargeCfgExist)
	}

	cfg.Name = req.Name
	cfg.Icon = req.Icon
	cfg.Diamond = req.Diamond
	cfg.ExtraDiamond = req.ExtraDiamond
	cfg.Price = req.Price
	if req.Currency != "" {
		cfg.Currency = req.Currency
	}
	cfg.ProductId = req.ProductId
	cfg.Sort = req.Sort
	cfg.Description = req.Description

	if err := rechargecfgdao.Update(cfg); err != nil {
		return nil, err
	}
	reloadRechargeCfgCache()
	return &rechargecfgdto.UpdateRechargeCfgRes{Success: true}, nil
}

// Delete 删除充值配置
func Delete(_ context.Context, req *rechargecfgdto.DeleteRechargeCfgReq) (*rechargecfgdto.DeleteRechargeCfgRes, error) {
	if cfg := rechargecfgdao.GetById(req.ID); cfg == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	if err := rechargecfgdao.Delete(req.ID); err != nil {
		return nil, err
	}
	reloadRechargeCfgCache()
	return &rechargecfgdto.DeleteRechargeCfgRes{Success: true}, nil
}

// OnShelf 上架
func OnShelf(_ context.Context, req *rechargecfgdto.OnShelfRechargeCfgReq) (*rechargecfgdto.OnShelfRechargeCfgRes, error) {
	cfg := rechargecfgdao.GetById(req.ID)
	if cfg == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	if cfg.Status != entity.RechargeCfgStatusOnShelf {
		if err := rechargecfgdao.UpdateStatus(req.ID, entity.RechargeCfgStatusOnShelf); err != nil {
			return nil, err
		}
		reloadRechargeCfgCache()
	}
	return &rechargecfgdto.OnShelfRechargeCfgRes{Success: true, Status: entity.RechargeCfgStatusOnShelf}, nil
}

// OffShelf 下架
func OffShelf(_ context.Context, req *rechargecfgdto.OffShelfRechargeCfgReq) (*rechargecfgdto.OffShelfRechargeCfgRes, error) {
	cfg := rechargecfgdao.GetById(req.ID)
	if cfg == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	if cfg.Status != entity.RechargeCfgStatusOffShelf {
		if err := rechargecfgdao.UpdateStatus(req.ID, entity.RechargeCfgStatusOffShelf); err != nil {
			return nil, err
		}
		reloadRechargeCfgCache()
	}
	return &rechargecfgdto.OffShelfRechargeCfgRes{Success: true, Status: entity.RechargeCfgStatusOffShelf}, nil
}

// ===== App =====

// GetAppList App端查询(仅返回已上架,走内存缓存)
// 缓存在服务启动时加载,CMS 端创建/修改/删除/上下架时重新从 DB 加载。
func GetAppList(_ context.Context, _ *rechargecfgdto.AppRechargeCfgListReq) (*rechargecfgdto.AppRechargeCfgListRes, error) {
	return &rechargecfgdto.AppRechargeCfgListRes{List: getRechargeCfgCache()}, nil
}
