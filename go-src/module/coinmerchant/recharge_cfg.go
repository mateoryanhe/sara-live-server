package coinmerchant

import (
	"context"
	"strconv"
	"strings"

	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dto/coinmerchantrechargecfgdto"
	"xr-game-server/entity/recharge"
	"xr-game-server/errercode"
)

// ListCoinMerchantRechargeCfgs CMS分页列表
func ListCoinMerchantRechargeCfgs(_ context.Context, req *coinmerchantrechargecfgdto.CoinMerchantRechargeCfgListReq) (*httpserver.CMSQueryResp, error) {
	total, list := cfgdao.GetCoinMerchantRechargeCfgList(req)
	return httpserver.NewCMSQueryResp(total, list), nil
}

// CreateCoinMerchantRechargeCfg 新建(默认下架)
func CreateCoinMerchantRechargeCfg(_ context.Context, req *coinmerchantrechargecfgdto.CreateCoinMerchantRechargeCfgReq) (*coinmerchantrechargecfgdto.CreateCoinMerchantRechargeCfgRes, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" || req.Price <= 0 || req.Gold == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if cfgdao.GetCoinMerchantRechargeCfgByName(name) != nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgExist)
	}
	cfg := &entity.CoinMerchantRechargeCfg{
		Name:   name,
		Price:  req.Price,
		Gold:   req.Gold,
		Status: entity.CoinMerchantRechargeCfgStatusOffShelf,
	}
	if err := cfgdao.CreateCoinMerchantRechargeCfg(cfg); err != nil {
		return nil, err
	}
	reloadCoinMerchantRechargeCfgCache()
	return &coinmerchantrechargecfgdto.CreateCoinMerchantRechargeCfgRes{ID: strconv.FormatUint(cfg.ID, 10)}, nil
}

// UpdateCoinMerchantRechargeCfg 修改(写库并刷缓存,不改上下架)
func UpdateCoinMerchantRechargeCfg(_ context.Context, req *coinmerchantrechargecfgdto.UpdateCoinMerchantRechargeCfgReq) (*coinmerchantrechargecfgdto.UpdateCoinMerchantRechargeCfgRes, error) {
	cfg := cfgdao.GetCoinMerchantRechargeCfgById(req.ID)
	if cfg == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	name := strings.TrimSpace(req.Name)
	if name == "" || req.Price <= 0 || req.Gold == 0 {
		return nil, errercode.CreateCode(errercode.InvalidParam)
	}
	if existing := cfgdao.GetCoinMerchantRechargeCfgByName(name); existing != nil && existing.ID != req.ID {
		return nil, errercode.CreateCode(errercode.RechargeCfgExist)
	}
	cfg.Name = name
	cfg.Price = req.Price
	cfg.Gold = req.Gold
	if err := cfgdao.UpdateCoinMerchantRechargeCfg(cfg); err != nil {
		return nil, err
	}
	reloadCoinMerchantRechargeCfgCache()
	return &coinmerchantrechargecfgdto.UpdateCoinMerchantRechargeCfgRes{Success: true}, nil
}

// DeleteCoinMerchantRechargeCfg 删除并刷缓存
func DeleteCoinMerchantRechargeCfg(_ context.Context, req *coinmerchantrechargecfgdto.DeleteCoinMerchantRechargeCfgReq) (*coinmerchantrechargecfgdto.DeleteCoinMerchantRechargeCfgRes, error) {
	if cfgdao.GetCoinMerchantRechargeCfgById(req.ID) == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	if err := cfgdao.DeleteCoinMerchantRechargeCfg(req.ID); err != nil {
		return nil, err
	}
	reloadCoinMerchantRechargeCfgCache()
	return &coinmerchantrechargecfgdto.DeleteCoinMerchantRechargeCfgRes{Success: true}, nil
}

// OnShelfCoinMerchantRechargeCfg 上架并刷缓存
func OnShelfCoinMerchantRechargeCfg(_ context.Context, req *coinmerchantrechargecfgdto.OnShelfCoinMerchantRechargeCfgReq) (*coinmerchantrechargecfgdto.OnShelfCoinMerchantRechargeCfgRes, error) {
	cfg := cfgdao.GetCoinMerchantRechargeCfgById(req.ID)
	if cfg == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	if cfg.Status != entity.CoinMerchantRechargeCfgStatusOnShelf {
		if err := cfgdao.UpdateCoinMerchantRechargeCfgStatus(req.ID, entity.CoinMerchantRechargeCfgStatusOnShelf); err != nil {
			return nil, err
		}
		reloadCoinMerchantRechargeCfgCache()
	}
	return &coinmerchantrechargecfgdto.OnShelfCoinMerchantRechargeCfgRes{
		Success: true,
		Status:  entity.CoinMerchantRechargeCfgStatusOnShelf,
	}, nil
}

// OffShelfCoinMerchantRechargeCfg 下架并刷缓存
func OffShelfCoinMerchantRechargeCfg(_ context.Context, req *coinmerchantrechargecfgdto.OffShelfCoinMerchantRechargeCfgReq) (*coinmerchantrechargecfgdto.OffShelfCoinMerchantRechargeCfgRes, error) {
	cfg := cfgdao.GetCoinMerchantRechargeCfgById(req.ID)
	if cfg == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	if cfg.Status != entity.CoinMerchantRechargeCfgStatusOffShelf {
		if err := cfgdao.UpdateCoinMerchantRechargeCfgStatus(req.ID, entity.CoinMerchantRechargeCfgStatusOffShelf); err != nil {
			return nil, err
		}
		reloadCoinMerchantRechargeCfgCache()
	}
	return &coinmerchantrechargecfgdto.OffShelfCoinMerchantRechargeCfgRes{
		Success: true,
		Status:  entity.CoinMerchantRechargeCfgStatusOffShelf,
	}, nil
}

// GetAppList App端查询(仅返回已上架,走内存缓存)
func GetAppList(_ context.Context, _ *coinmerchantrechargecfgdto.AppCoinMerchantRechargeCfgListReq) (*coinmerchantrechargecfgdto.AppCoinMerchantRechargeCfgListRes, error) {
	list := GetOnShelfCoinMerchantRechargeCfgs()
	if list == nil {
		list = emptyCoinMerchantRechargeCfgList
	}
	return &coinmerchantrechargecfgdto.AppCoinMerchantRechargeCfgListRes{List: list}, nil
}
