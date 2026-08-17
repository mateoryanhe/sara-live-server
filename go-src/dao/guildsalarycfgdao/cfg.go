package guildsalarycfgdao

import (
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/str"
	"xr-game-server/dto/guildsalarycfgdto"
	"xr-game-server/entity/live"
)

func GetById(id uint64) *entity.GuildSalaryCfg {
	if id == 0 {
		return nil
	}
	var row entity.GuildSalaryCfg
	if err := g.DB().Model(string(entity.TbGuildSalaryCfg)).Where("id = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func Create(row *entity.GuildSalaryCfg) error {
	_, err := g.DB().Model(string(entity.TbGuildSalaryCfg)).Save(row)
	return err
}

func Update(row *entity.GuildSalaryCfg) error {
	return Create(row)
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbGuildSalaryCfg)).WherePri(id).Delete()
	return err
}

func GetList(req *guildsalarycfgdto.GuildSalaryCfgListReq) (int, []*guildsalarycfgdto.GuildSalaryCfgItem) {
	sql := `select id, weekly_work_days, daily_live_duration_minutes, salary_amount, sort, created_at, updated_at
            from guild_salary_cfgs
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*guildsalarycfgdto.GuildSalaryCfgItem, 0)

	sql += ` order by salary_amount desc, id desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}

// ListAllOrderBySalaryDesc 全部工会薪资分档(薪资降序,用于结算匹配最高档)
func ListAllOrderBySalaryDesc() []*entity.GuildSalaryCfg {
	rows := make([]*entity.GuildSalaryCfg, 0)
	_ = g.DB().Model(string(entity.TbGuildSalaryCfg)).
		Order("salary_amount desc, id desc").
		Scan(&rows)
	return rows
}
