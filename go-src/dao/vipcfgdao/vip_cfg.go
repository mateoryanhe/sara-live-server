package vipcfgdao

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"strconv"
	"xr-game-server/core/str"
	"xr-game-server/dto/vipcfgdto"
	"xr-game-server/entity"
)

func GetById(id uint64) *entity.VipCfg {
	var row entity.VipCfg
	if err := g.DB().Model(string(entity.TbVipCfg)).Where("id = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func GetByLevel(level uint32) *entity.VipCfg {
	var row entity.VipCfg
	if err := g.DB().Model(string(entity.TbVipCfg)).Where("level = ?", level).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func Create(row *entity.VipCfg) error {
	_, err := g.DB().Model(string(entity.TbVipCfg)).Save(row)
	return err
}

func Update(row *entity.VipCfg) error {
	return Create(row)
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbVipCfg)).WherePri(id).Delete()
	return err
}

// GetAll 查询全部VIP配置(按等级升序,供内存缓存加载)
func GetAll() []*entity.VipCfg {
	var rows []*entity.VipCfg
	_ = g.DB().Model(string(entity.TbVipCfg)).Order("level asc").Scan(&rows)
	return rows
}

func GetList(req *vipcfgdto.VipCfgListReq) (int, []*vipcfgdto.VipCfgListRes) {
	sql := `select id, level, level_name, status, upgrade_recharge_limit, min_withdraw_amount,
                   max_withdraw_amount, fee, created_at, updated_at
            from vip_cfgs
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*vipcfgdto.VipCfgListRes, 0)

	if req.LevelName != "" {
		sql += ` and level_name LIKE ?`
		param = append(param, fmt.Sprintf("%%%s%%", req.LevelName))
	}
	switch req.StatusFilter {
	case 1:
		sql += ` and status = ?`
		param = append(param, entity.VipCfgStatusDisabled)
	case 2:
		sql += ` and status = ?`
		param = append(param, entity.VipCfgStatusEnabled)
	}

	sql += ` order by level asc, created_at desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
