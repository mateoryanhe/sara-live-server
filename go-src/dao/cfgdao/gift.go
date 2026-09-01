package cfgdao

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"strconv"
	"xr-game-server/core/str"
	"xr-game-server/dto/giftdto"
	"xr-game-server/entity/live"
)

func GetGiftById(id uint64) *entity.LiveGift {
	var gift entity.LiveGift
	err := g.DB().Model(string(entity.TbLiveGift)).Where("id = ?", id).Scan(&gift)
	if err != nil {
		return nil
	}
	return &gift
}

func GetGiftByName(name string) *entity.LiveGift {
	var gift entity.LiveGift
	err := g.DB().Model(string(entity.TbLiveGift)).Where("name = ?", name).Scan(&gift)
	if err != nil {
		return nil
	}
	return &gift
}

func CreateGift(gift *entity.LiveGift) error {
	_, err := g.DB().Model(string(entity.TbLiveGift)).Save(gift)
	return err
}

func GetAllGiftsForAssetPreview() []*entity.LiveGift {
	var rows []*entity.LiveGift
	_ = g.DB().Model(string(entity.TbLiveGift)).
		Fields("id, name, icon, animation, status, sort").
		Order("sort desc, id asc").
		Scan(&rows)
	return rows
}

func GetGiftsByIDs(ids []uint64) []*entity.LiveGift {
	if len(ids) == 0 {
		return nil
	}
	var rows []*entity.LiveGift
	_ = g.DB().Model(string(entity.TbLiveGift)).WhereIn("id", ids).Scan(&rows)
	return rows
}

func UpdateGift(gift *entity.LiveGift) error {
	return CreateGift(gift)
}

func DeleteGift(id uint64) error {
	_, err := g.DB().Model(string(entity.TbLiveGift)).WherePri(id).Delete()
	return err
}

func UpdateGiftStatus(id uint64, status uint8) error {
	_, err := g.DB().Model(string(entity.TbLiveGift)).
		WherePri(id).
		Data(g.Map{"status": status}).
		Update()
	return err
}

func GetOnShelfGifts() []*entity.LiveGift {
	ret := make([]*entity.LiveGift, 0)
	now := time.Now()
	err := g.DB().Model(string(entity.TbLiveGift)).
		Where("status = ? AND (published_at IS NULL OR published_at <= ?)", entity.LiveGiftStatusOnShelf, now).
		Order("sort desc, published_at desc, created_at desc").
		Scan(&ret)
	if err != nil {
		return nil
	}
	return ret
}

func GetGiftList(req *giftdto.GiftListReq) (int, []*giftdto.GiftListRes) {
	sql := `select id, name, name_en, name_es, name_pt, name_hi, name_id, icon, animation, price, category, sort, status, published_at, description, created_at, updated_at
            from live_gifts
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*giftdto.GiftListRes, 0)

	if req.Name != "" {
		sql += ` and name LIKE ?`
		param = append(param, fmt.Sprintf("%%%s%%", req.Name))
	}
	if req.Category != "" {
		sql += ` and category = ?`
		param = append(param, req.Category)
	}
	switch req.StatusFilter {
	case 1:
		sql += ` and status = ?`
		param = append(param, entity.LiveGiftStatusOffShelf)
	case 2:
		sql += ` and status = ?`
		param = append(param, entity.LiveGiftStatusOnShelf)
	}

	sql += ` order by sort desc, created_at desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageOffset())
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
