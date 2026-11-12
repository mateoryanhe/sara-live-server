package anchorvideopricetierdao

import (
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/str"
	"xr-game-server/dto/anchorvideopricetierdto"
	entity "xr-game-server/entity/live"
)

func GetById(id uint64) *entity.AnchorVideoPriceTier {
	var row entity.AnchorVideoPriceTier
	if err := g.DB().Model(string(entity.TbAnchorVideoPriceTier)).Where("id = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func GetAll() []*entity.AnchorVideoPriceTier {
	ret := make([]*entity.AnchorVideoPriceTier, 0)
	_ = g.DB().Model(string(entity.TbAnchorVideoPriceTier)).
		Order("price asc, created_at desc").
		Scan(&ret)
	return ret
}

func Create(row *entity.AnchorVideoPriceTier) error {
	_, err := g.DB().Model(string(entity.TbAnchorVideoPriceTier)).Save(row)
	return err
}

func Update(row *entity.AnchorVideoPriceTier) error {
	return Create(row)
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbAnchorVideoPriceTier)).WherePri(id).Delete()
	return err
}

func UpdateStatus(id uint64, status uint8) error {
	_, err := g.DB().Model(string(entity.TbAnchorVideoPriceTier)).
		WherePri(id).
		Data(g.Map{"status": status}).
		Update()
	return err
}

func GetList(req *anchorvideopricetierdto.PriceTierListReq) (int, []*anchorvideopricetierdto.PriceTierListItem) {
	sql := `select id, price, status, created_at, updated_at
            from anchor_video_price_tiers
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*anchorvideopricetierdto.PriceTierListItem, 0)

	switch req.StatusFilter {
	case 1:
		sql += ` and status = ?`
		param = append(param, entity.AnchorVideoPriceTierStatusOffShelf)
	case 2:
		sql += ` and status = ?`
		param = append(param, entity.AnchorVideoPriceTierStatusOnShelf)
	}

	sql += ` order by price asc, created_at desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageOffset())
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
