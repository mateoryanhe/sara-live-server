package recharge

import (
	"context"
	"strconv"
	"strings"
	"xr-game-server/core/httpserver"
	"xr-game-server/dao/cfgdao"
	"xr-game-server/dao/userinfodao"
	"xr-game-server/dto/rechargecfgdto"
	"xr-game-server/entity/recharge"
	"xr-game-server/errercode"
	"xr-game-server/module/activity"
	"xr-game-server/module/upload"
)

// ===== CMS =====

// GetList CMS分页查询(全部状态)
func GetList(_ context.Context, req *rechargecfgdto.RechargeCfgListReq) (*httpserver.CMSQueryResp, error) {
	total, list := cfgdao.GetRechargeCfgList(req)
	for _, row := range list {
		row.IconName = row.Icon
		row.Icon = upload.GetUrlByName(row.IconName)
	}
	return &httpserver.CMSQueryResp{Total: total, Data: list}, nil
}

// Create 创建充值配置(默认下架,需手动上架)
func Create(_ context.Context, req *rechargecfgdto.CreateRechargeCfgReq) (*rechargecfgdto.CreateRechargeCfgRes, error) {
	if existing := cfgdao.GetRechargeCfgByNameAndType(req.Name, req.CfgType); existing != nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgExist)
	}
	if err := validateRechargeCfgProductId(req.ProductId, req.CfgType, 0); err != nil {
		return nil, err
	}
	cfg := &entity.RechargeCfg{
		Name:        req.Name,
		CfgType:     req.CfgType,
		Icon:        req.Icon,
		Gold:        req.Gold,
		Price:       req.Price,
		Currency:    entity.RechargeCfgCurrencyUSD,
		ProductId:   req.ProductId,
		Status:      entity.RechargeCfgStatusOffShelf,
		Description: req.Description,
	}
	if err := cfgdao.CreateRechargeCfg(cfg); err != nil {
		return nil, err
	}
	reloadRechargeCfgCache()
	return &rechargecfgdto.CreateRechargeCfgRes{ID: strconv.FormatUint(cfg.ID, 10)}, nil
}

// Update 修改充值配置(不修改上下架状态)
func Update(_ context.Context, req *rechargecfgdto.UpdateRechargeCfgReq) (*rechargecfgdto.UpdateRechargeCfgRes, error) {
	cfg := cfgdao.GetRechargeCfgById(req.ID)
	if cfg == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	if existing := cfgdao.GetRechargeCfgByNameAndType(req.Name, req.CfgType); existing != nil && existing.ID != req.ID {
		return nil, errercode.CreateCode(errercode.RechargeCfgExist)
	}
	if err := validateRechargeCfgProductId(req.ProductId, req.CfgType, req.ID); err != nil {
		return nil, err
	}

	cfg.Name = req.Name
	cfg.CfgType = req.CfgType
	cfg.Icon = req.Icon
	cfg.Gold = req.Gold
	cfg.Price = req.Price
	cfg.Currency = entity.RechargeCfgCurrencyUSD
	cfg.ProductId = req.ProductId
	cfg.Description = req.Description

	if err := cfgdao.UpdateRechargeCfg(cfg); err != nil {
		return nil, err
	}
	reloadRechargeCfgCache()
	return &rechargecfgdto.UpdateRechargeCfgRes{Success: true}, nil
}

// Delete 删除充值配置
func Delete(_ context.Context, req *rechargecfgdto.DeleteRechargeCfgReq) (*rechargecfgdto.DeleteRechargeCfgRes, error) {
	if cfg := cfgdao.GetRechargeCfgById(req.ID); cfg == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	if err := cfgdao.DeleteRechargeCfg(req.ID); err != nil {
		return nil, err
	}
	reloadRechargeCfgCache()
	return &rechargecfgdto.DeleteRechargeCfgRes{Success: true}, nil
}

// OnShelf 上架
func OnShelf(_ context.Context, req *rechargecfgdto.OnShelfRechargeCfgReq) (*rechargecfgdto.OnShelfRechargeCfgRes, error) {
	cfg := cfgdao.GetRechargeCfgById(req.ID)
	if cfg == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	if cfg.Status != entity.RechargeCfgStatusOnShelf {
		if err := cfgdao.UpdateRechargeCfgStatus(req.ID, entity.RechargeCfgStatusOnShelf); err != nil {
			return nil, err
		}
		reloadRechargeCfgCache()
	}
	return &rechargecfgdto.OnShelfRechargeCfgRes{Success: true, Status: entity.RechargeCfgStatusOnShelf}, nil
}

// OffShelf 下架
func OffShelf(_ context.Context, req *rechargecfgdto.OffShelfRechargeCfgReq) (*rechargecfgdto.OffShelfRechargeCfgRes, error) {
	cfg := cfgdao.GetRechargeCfgById(req.ID)
	if cfg == nil {
		return nil, errercode.CreateCode(errercode.RechargeCfgNonExist)
	}
	if cfg.Status != entity.RechargeCfgStatusOffShelf {
		if err := cfgdao.UpdateRechargeCfgStatus(req.ID, entity.RechargeCfgStatusOffShelf); err != nil {
			return nil, err
		}
		reloadRechargeCfgCache()
	}
	return &rechargecfgdto.OffShelfRechargeCfgRes{Success: true, Status: entity.RechargeCfgStatusOffShelf}, nil
}

// ===== App =====

// GetAppList App端查询(仅返回已上架,走内存缓存)
// 缓存在服务启动时加载,CMS 端创建/修改/删除/上下架时重新从 DB 加载。
func GetAppList(ctx context.Context, req *rechargecfgdto.AppRechargeCfgListReq) (*rechargecfgdto.AppRechargeCfgListRes, error) {
	_ = req
	return &rechargecfgdto.AppRechargeCfgListRes{
		List: buildAppRechargeCfgList(httpserver.GetAuthId(ctx)),
	}, nil
}

// GetAppListByUserId App端按上报用户ID查询(无需鉴权,仅返回已上架,走内存缓存)
func GetAppListByUserId(_ context.Context, req *rechargecfgdto.AppRechargeCfgListByUserIdReq) (*rechargecfgdto.AppRechargeCfgListRes, error) {
	userId, err := strconv.ParseUint(strings.TrimSpace(req.UserId), 10, 64)
	if err != nil || userId == 0 {
		return nil, errercode.CreateCode(errercode.EmptyUserId)
	}
	if userinfodao.GetUserInfoByUserId(userId) == nil {
		return nil, errercode.CreateCode(errercode.SysError)
	}
	return &rechargecfgdto.AppRechargeCfgListRes{
		List: buildAppRechargeCfgList(userId),
	}, nil
}

func buildAppRechargeCfgList(userId uint64) []*rechargecfgdto.AppRechargeCfgItem {
	cfgRatio := activity.ConfiguredFirstRechargeRatio()
	all := getRechargeCfgCache()
	list := make([]*rechargecfgdto.AppRechargeCfgItem, 0, len(all))
	for _, item := range all {
		if item == nil {
			continue
		}
		copyItem := *item
		if userId > 0 {
			copyItem.FirstRecharge = userinfodao.IsRechargeCfgFirstRecharge(userId, item.ID)
		} else {
			copyItem.FirstRecharge = true
		}
		if copyItem.FirstRecharge {
			copyItem.FirstRechargeRatio = cfgRatio
		} else {
			copyItem.FirstRechargeRatio = 0
		}
		list = append(list, &copyItem)
	}
	return list
}

func validateRechargeCfgProductId(productId string, cfgType uint8, excludeID uint64) error {
	productId = strings.TrimSpace(productId)
	if productId == "" {
		return nil
	}
	switch cfgType {
	case entity.RechargeCfgTypeGoogle, entity.RechargeCfgTypeIOS:
	default:
		return nil
	}
	existing := cfgdao.GetRechargeCfgByProductIdAndType(productId, cfgType)
	if existing != nil && existing.ID != excludeID {
		return errercode.CreateCode(errercode.RechargeCfgProductIdExist)
	}
	return nil
}

