package shortvideodao

import (
	"strconv"
	"xr-game-server/core/str"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity/shortvideo"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func GetPriceTierById(id uint64) *entity.ShortVideoPriceTier {
	var row entity.ShortVideoPriceTier
	if err := g.DB().Model(string(entity.TbShortVideoPriceTier)).Where("id = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func GetAllPriceTiers() []*entity.ShortVideoPriceTier {
	ret := make([]*entity.ShortVideoPriceTier, 0)
	_ = g.DB().Model(string(entity.TbShortVideoPriceTier)).
		Order("price asc, created_at desc").
		Scan(&ret)
	return ret
}

func CreatePriceTier(row *entity.ShortVideoPriceTier) error {
	_, err := g.DB().Model(string(entity.TbShortVideoPriceTier)).Save(row)
	return err
}

func UpdatePriceTier(row *entity.ShortVideoPriceTier) error {
	return CreatePriceTier(row)
}

func DeletePriceTier(id uint64) error {
	_, err := g.DB().Model(string(entity.TbShortVideoPriceTier)).WherePri(id).Delete()
	return err
}

func UpdatePriceTierStatus(id uint64, status uint8) error {
	_, err := g.DB().Model(string(entity.TbShortVideoPriceTier)).
		WherePri(id).
		Data(g.Map{"status": status}).
		Update()
	return err
}

func GetPriceTierList(req *shortvideodto.ShortVideoPriceTierListReq) (int, []*shortvideodto.ShortVideoPriceTierListRes) {
	sql := `select id, price, status, created_at, updated_at
            from short_video_price_tiers
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*shortvideodto.ShortVideoPriceTierListRes, 0)

	switch req.StatusFilter {
	case 1:
		sql += ` and status = ?`
		param = append(param, entity.ShortVideoPriceTierStatusOffShelf)
	case 2:
		sql += ` and status = ?`
		param = append(param, entity.ShortVideoPriceTierStatusOnShelf)
	}

	sql += ` order by price asc, created_at desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageOffset())
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
