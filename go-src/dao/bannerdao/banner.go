package bannerdao

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"strconv"
	"xr-game-server/core/str"
	"xr-game-server/dto/bannerdto"
	"xr-game-server/entity"
)

func GetById(id uint64) *entity.HomeBanner {
	var row entity.HomeBanner
	if err := g.DB().Model(string(entity.TbHomeBanner)).Where("id = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func GetByTitle(title string) *entity.HomeBanner {
	var row entity.HomeBanner
	if err := g.DB().Model(string(entity.TbHomeBanner)).Where("title = ?", title).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func GetAll() []*entity.HomeBanner {
	ret := make([]*entity.HomeBanner, 0)
	_ = g.DB().Model(string(entity.TbHomeBanner)).
		Order("sort desc, created_at desc").
		Scan(&ret)
	return ret
}

func Create(row *entity.HomeBanner) error {
	_, err := g.DB().Model(string(entity.TbHomeBanner)).Save(row)
	return err
}

func Update(row *entity.HomeBanner) error {
	return Create(row)
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbHomeBanner)).WherePri(id).Delete()
	return err
}

func UpdateStatus(id uint64, status uint8) error {
	_, err := g.DB().Model(string(entity.TbHomeBanner)).
		WherePri(id).
		Data(g.Map{"status": status}).
		Update()
	return err
}

func GetBannerList(req *bannerdto.BannerListReq) (int, []*bannerdto.BannerListRes) {
	sql := `select id, title, image, link, scene, direction, sort, status, created_at, updated_at
            from home_banners
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*bannerdto.BannerListRes, 0)

	if req.Title != "" {
		sql += ` and title LIKE ?`
		param = append(param, fmt.Sprintf("%%%s%%", req.Title))
	}
	switch req.SceneFilter {
	case int(entity.HomeBannerSceneHome):
		sql += ` and scene = ?`
		param = append(param, entity.HomeBannerSceneHome)
	case int(entity.HomeBannerSceneLiveRoom):
		sql += ` and scene = ?`
		param = append(param, entity.HomeBannerSceneLiveRoom)
	}
	switch req.StatusFilter {
	case 1:
		sql += ` and status = ?`
		param = append(param, entity.HomeBannerStatusOffShelf)
	case 2:
		sql += ` and status = ?`
		param = append(param, entity.HomeBannerStatusOnShelf)
	}

	sql += ` order by sort desc, created_at desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
