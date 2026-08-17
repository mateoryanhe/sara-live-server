package anchorsalarycfgdao

import (
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/str"
	"xr-game-server/dto/anchorsalarycfgdto"
	"xr-game-server/entity/live"
)

func GetById(id uint64) *entity.AnchorSalaryCfg {
	if id == 0 {
		return nil
	}
	var row entity.AnchorSalaryCfg
	if err := g.DB().Model(string(entity.TbAnchorSalaryCfg)).Where("id = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func Create(row *entity.AnchorSalaryCfg) error {
	_, err := g.DB().Model(string(entity.TbAnchorSalaryCfg)).Save(row)
	return err
}

func Update(row *entity.AnchorSalaryCfg) error {
	return Create(row)
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbAnchorSalaryCfg)).WherePri(id).Delete()
	return err
}

func GetList(req *anchorsalarycfgdto.AnchorSalaryCfgListReq) (int, []*anchorsalarycfgdto.AnchorSalaryCfgItem) {
	sql := `select id, daily_effective_live_count, weekly_effective_live_count, salary_amount, sort, created_at, updated_at
            from anchor_salary_cfgs
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*anchorsalarycfgdto.AnchorSalaryCfgItem, 0)

	sql += ` order by salary_amount desc, id desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
