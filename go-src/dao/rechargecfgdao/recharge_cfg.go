package rechargecfgdao

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"strconv"
	"xr-game-server/core/str"
	"xr-game-server/dto/rechargecfgdto"
	"xr-game-server/entity"
)

// GetById 按ID获取充值配置(直接查 DB,不走缓存)
func GetById(id uint64) *entity.RechargeCfg {
	var cfg entity.RechargeCfg
	err := g.DB().Model(string(entity.TbRechargeCfg)).Where("id = ?", id).Scan(&cfg)
	if err != nil {
		return nil
	}
	return &cfg
}

// GetByName 按名称获取(用于唯一性校验)
func GetByName(name string) *entity.RechargeCfg {
	var cfg entity.RechargeCfg
	err := g.DB().Model(string(entity.TbRechargeCfg)).Where("name = ?", name).Scan(&cfg)
	if err != nil {
		return nil
	}
	return &cfg
}

// Create 新建充值配置
func Create(cfg *entity.RechargeCfg) error {
	_, err := g.DB().Model(string(entity.TbRechargeCfg)).Save(cfg)
	return err
}

// Update 更新充值配置(整行 Save)
func Update(cfg *entity.RechargeCfg) error {
	return Create(cfg)
}

// Delete 删除充值配置
func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbRechargeCfg)).WherePri(id).Delete()
	return err
}

// UpdateStatus 仅更新上下架状态
func UpdateStatus(id uint64, status uint8) error {
	_, err := g.DB().Model(string(entity.TbRechargeCfg)).
		WherePri(id).
		Data(g.Map{"status": status}).
		Update()
	return err
}

// GetOnShelf 获取全部已上架配置(按 sort desc, price asc, created_at desc 排序)
// App 端拉取上架列表使用,直接走 DB(数据量小,变更不频繁,可承受每次查询)
func GetOnShelf() []*entity.RechargeCfg {
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

// GetList CMS 分页查询(支持名称模糊、上下架过滤)
func GetList(req *rechargecfgdto.RechargeCfgListReq) (int, []*rechargecfgdto.RechargeCfgListRes) {
	sql := `select id, name, icon, diamond, extra_diamond, price, currency, product_id,
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
	switch req.StatusFilter {
	case 1: // 只看下架
		sql += ` and status = ?`
		param = append(param, entity.RechargeCfgStatusOffShelf)
	case 2: // 只看上架
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
