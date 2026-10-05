package anchornosalarysharecfgdao

import (
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/dto/anchornosalarysharecfgdto"
	"xr-game-server/entity/live"
)

func CountAll() int {
	total, _ := g.DB().Model(string(entity.TbAnchorNoSalaryShareCfg)).Count()
	return total
}

func GetByID(id uint64) *entity.AnchorNoSalaryShareCfg {
	if id == 0 {
		return nil
	}
	var row entity.AnchorNoSalaryShareCfg
	if err := g.DB().Model(string(entity.TbAnchorNoSalaryShareCfg)).Where("id = ?", id).Scan(&row); err != nil || row.ID == 0 {
		return nil
	}
	return &row
}

func Create(row *entity.AnchorNoSalaryShareCfg) error {
	_, err := g.DB().Model(string(entity.TbAnchorNoSalaryShareCfg)).Save(row)
	return err
}

func Update(row *entity.AnchorNoSalaryShareCfg) error {
	return Create(row)
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbAnchorNoSalaryShareCfg)).Where("id = ?", id).Delete()
	return err
}

func Exists(level uint32, excludeID uint64) bool {
	model := g.DB().Model(string(entity.TbAnchorNoSalaryShareCfg)).Where("level = ?", level)
	if excludeID > 0 {
		model = model.Where("id <> ?", excludeID)
	}
	total, err := model.Count()
	return err == nil && total > 0
}

// ListAllOrderByLevelAsc 结算时按等级顺序读取；当前等级的流水值是升级到下一等级的边界。
func ListAllOrderByLevelAsc() []*entity.AnchorNoSalaryShareCfg {
	rows := make([]*entity.AnchorNoSalaryShareCfg, 0)
	_ = g.DB().Model(string(entity.TbAnchorNoSalaryShareCfg)).
		Order("level asc, id asc").Scan(&rows)
	return rows
}

func GetList(req *anchornosalarysharecfgdto.AnchorNoSalaryShareCfgListReq) (int, []*anchornosalarysharecfgdto.AnchorNoSalaryShareCfgItem) {
	sql := `select id, level, social_total_diamond_revenue, anchor_social_share_percent, guild_social_share_percent, created_at, updated_at
	            from anchor_no_salary_share_cfgs
	            order by level asc, id asc`
	ctx := gctx.New()
	total, _ := g.DB().Model(string(entity.TbAnchorNoSalaryShareCfg)).Count()
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageOffset())
	ret := make([]*anchornosalarysharecfgdto.AnchorNoSalaryShareCfgItem, 0)
	g.DB().GetScan(ctx, &ret, sql)
	return total, ret
}
