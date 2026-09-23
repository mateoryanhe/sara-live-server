package anchorgamesharecfgdao

import (
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/str"
	"xr-game-server/dto/anchorgamesharecfgdto"
	"xr-game-server/entity/live"
)

func Exists(salaryType, level uint32, excludeID uint64) bool {
	model := g.DB().Model(string(entity.TbAnchorGameShareCfg)).
		Where("salary_type = ? and level = ?", salaryType, level)
	if excludeID > 0 {
		model = model.Where("id <> ?", excludeID)
	}
	total, _ := model.Count()
	return total > 0
}

func GetByID(id uint64) *entity.AnchorGameShareCfg {
	if id == 0 {
		return nil
	}
	var row entity.AnchorGameShareCfg
	if err := g.DB().Model(string(entity.TbAnchorGameShareCfg)).Where("id = ?", id).Scan(&row); err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

func Create(row *entity.AnchorGameShareCfg) error {
	_, err := g.DB().Model(string(entity.TbAnchorGameShareCfg)).Save(row)
	return err
}

func Update(row *entity.AnchorGameShareCfg) error {
	return Create(row)
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbAnchorGameShareCfg)).WherePri(id).Delete()
	return err
}

// ListAllBySalaryTypeOrderByThresholdDesc 结算时按金币流水门槛由高到低匹配最高档。
func ListAllBySalaryTypeOrderByThresholdDesc(salaryType uint32) []*entity.AnchorGameShareCfg {
	rows := make([]*entity.AnchorGameShareCfg, 0)
	_ = g.DB().Model(string(entity.TbAnchorGameShareCfg)).
		Where("salary_type = ?", salaryType).
		Order("game_total_gold_revenue desc, level desc, id desc").Scan(&rows)
	return rows
}

func GetList(req *anchorgamesharecfgdto.AnchorGameShareCfgListReq) (int, []*anchorgamesharecfgdto.AnchorGameShareCfgItem) {
	sql := `select id, salary_type, level, game_total_gold_revenue, anchor_game_share_percent, guild_game_share_percent, created_at, updated_at
            from anchor_game_share_cfgs
            where salary_type = ` + strconv.FormatUint(uint64(req.SalaryType), 10) + `
            order by level asc, id asc`
	ctx := gctx.New()
	countSQL := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSQL, nil)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageOffset())
	ret := make([]*anchorgamesharecfgdto.AnchorGameShareCfgItem, 0)
	g.DB().GetScan(ctx, &ret, sql, nil)
	return total, ret
}
