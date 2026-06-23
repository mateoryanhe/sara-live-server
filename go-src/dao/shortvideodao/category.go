package shortvideodao

import (
	"strconv"
	"xr-game-server/core/str"
	"xr-game-server/dto/shortvideodto"
	"xr-game-server/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func GetCategoryById(id uint64) *entity.ShortVideoCategory {
	var row entity.ShortVideoCategory
	if err := g.DB().Model(string(entity.TbShortVideoCategory)).Where("id = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func GetCategoryByName(name string) *entity.ShortVideoCategory {
	if name == "" {
		return nil
	}
	var row entity.ShortVideoCategory
	if err := g.DB().Model(string(entity.TbShortVideoCategory)).Where("name = ?", name).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func GetAllCategories() []*entity.ShortVideoCategory {
	ret := make([]*entity.ShortVideoCategory, 0)
	_ = g.DB().Model(string(entity.TbShortVideoCategory)).
		Order("sort desc, created_at desc").
		Scan(&ret)
	return ret
}

func CreateCategory(row *entity.ShortVideoCategory) error {
	_, err := g.DB().Model(string(entity.TbShortVideoCategory)).Save(row)
	return err
}

func UpdateCategory(row *entity.ShortVideoCategory) error {
	return CreateCategory(row)
}

func DeleteCategory(id uint64) error {
	_, err := g.DB().Model(string(entity.TbShortVideoCategory)).WherePri(id).Delete()
	return err
}

func GetCategoryList(req *shortvideodto.ShortVideoCategoryListReq) (int, []*shortvideodto.ShortVideoCategoryListRes) {
	sql := `select id, name, sort, created_at, updated_at
            from short_video_categories
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*shortvideodto.ShortVideoCategoryListRes, 0)

	sql += ` order by sort desc, created_at desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
