package cfgdao

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"strconv"
	"strings"
	"xr-game-server/core/str"
	"xr-game-server/dto/rechargecfgdto"
	"xr-game-server/entity"
)

func GetRechargeCfgById(id uint64) *entity.RechargeCfg {
	var cfg entity.RechargeCfg
	err := g.DB().Model(string(entity.TbRechargeCfg)).Where("id = ?", id).Scan(&cfg)
	if err != nil {
		return nil
	}
	return &cfg
}

func GetRechargeCfgByNameTypeAndPackage(name string, cfgType uint8, packageName string) *entity.RechargeCfg {
	var cfg entity.RechargeCfg
	err := g.DB().Model(string(entity.TbRechargeCfg)).
		Where("name = ? AND cfg_type = ? AND package_name = ?", name, cfgType, packageName).
		Scan(&cfg)
	if err != nil {
		return nil
	}
	if cfg.ID == 0 {
		return nil
	}
	return &cfg
}

func GetRechargeCfgByProductIdTypeAndPackage(productId string, cfgType uint8, packageName string) *entity.RechargeCfg {
	productId = strings.TrimSpace(productId)
	packageName = strings.TrimSpace(packageName)
	if productId == "" || packageName == "" {
		return nil
	}
	var cfg entity.RechargeCfg
	err := g.DB().Model(string(entity.TbRechargeCfg)).
		Where("product_id = ? AND cfg_type = ? AND package_name = ?", productId, cfgType, packageName).
		Scan(&cfg)
	if err != nil {
		return nil
	}
	if cfg.ID == 0 {
		return nil
	}
	return &cfg
}

func CreateRechargeCfg(cfg *entity.RechargeCfg) error {
	_, err := g.DB().Model(string(entity.TbRechargeCfg)).Save(cfg)
	return err
}

func UpdateRechargeCfg(cfg *entity.RechargeCfg) error {
	return CreateRechargeCfg(cfg)
}

func DeleteRechargeCfg(id uint64) error {
	_, err := g.DB().Model(string(entity.TbRechargeCfg)).WherePri(id).Delete()
	return err
}

func UpdateRechargeCfgStatus(id uint64, status uint8) error {
	_, err := g.DB().Model(string(entity.TbRechargeCfg)).
		WherePri(id).
		Data(g.Map{"status": status}).
		Update()
	return err
}

func GetOnShelfRechargeCfg() []*entity.RechargeCfg {
	ret := make([]*entity.RechargeCfg, 0)
	err := g.DB().Model(string(entity.TbRechargeCfg)).
		Where("status = ?", entity.RechargeCfgStatusOnShelf).
		Order("sort desc, price asc, created_at desc").
		Scan(&ret)
	if err != nil {
		return nil
	}
	return ret
}

func GetRechargeCfgList(req *rechargecfgdto.RechargeCfgListReq) (int, []*rechargecfgdto.RechargeCfgListRes) {
	sql := `select id, name, package_name, cfg_type, icon, gold, extra_gold, price, currency, product_id,
                   sort, status, description, created_at, updated_at
            from recharge_cfgs
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*rechargecfgdto.RechargeCfgListRes, 0)

	if req.Name != "" {
		sql += ` and name LIKE ?`
		param = append(param, fmt.Sprintf("%%%s%%", req.Name))
	}
	if req.PackageName != "" {
		sql += ` and package_name LIKE ?`
		param = append(param, fmt.Sprintf("%%%s%%", req.PackageName))
	}
	switch req.TypeFilter {
	case 1, 2, 3:
		sql += ` and cfg_type = ?`
		param = append(param, req.TypeFilter)
	}
	switch req.StatusFilter {
	case 1:
		sql += ` and status = ?`
		param = append(param, entity.RechargeCfgStatusOffShelf)
	case 2:
		sql += ` and status = ?`
		param = append(param, entity.RechargeCfgStatusOnShelf)
	}

	sql += ` order by sort desc, price asc, created_at desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
