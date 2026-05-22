package giftdao

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"strconv"
	"xr-game-server/core/str"
	"xr-game-server/dto/giftdto"
	"xr-game-server/entity"
)

// GetGiftById 根据ID获取礼物
func GetGiftById(id uint64) *entity.LiveGift {
	var gift entity.LiveGift
	err := g.DB().Model(string(entity.TbLiveGift)).Where("id = ?", id).Scan(&gift)
	if err != nil {
		return nil
	}
	return &gift
}

// GetGiftByName 根据名称获取礼物
func GetGiftByName(name string) *entity.LiveGift {
	var gift entity.LiveGift
	err := g.DB().Model(string(entity.TbLiveGift)).Where("name = ?", name).Scan(&gift)
	if err != nil {
		return nil
	}
	return &gift
}

// CreateGift 创建礼物
func CreateGift(gift *entity.LiveGift) error {
	_, err := g.DB().Model(string(entity.TbLiveGift)).Save(gift)
	if err == nil {
		return nil
	}
	return err
}

// UpdateGift 更新礼物(整行 Save)
func UpdateGift(gift *entity.LiveGift) error {
	return CreateGift(gift)
}

// DeleteGift 删除礼物
func DeleteGift(id uint64) error {
	_, err := g.DB().Model(string(entity.TbLiveGift)).WherePri(id).Delete()
	if err == nil {
		return nil
	}
	return err
}

// UpdateGiftStatus 仅更新上下架状态
func UpdateGiftStatus(id uint64, status uint8) error {
	_, err := g.DB().Model(string(entity.TbLiveGift)).
		WherePri(id).
		Data(g.Map{"status": status}).
		Update()
	return err
}

// GetOnShelfGifts 获取所有已上架且已到发布时间的礼物
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

// GetGiftList 分页获取礼物列表
func GetGiftList(req *giftdto.GiftListReq) (int, []*giftdto.GiftListRes) {
	sql := `select id, name, icon, animation, price, category, sort, status, published_at, description, created_at, updated_at
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
	case 1: // 只看下架
		sql += ` and status = ?`
		param = append(param, entity.LiveGiftStatusOffShelf)
	case 2: // 只看上架
		sql += ` and status = ?`
		param = append(param, entity.LiveGiftStatusOnShelf)
	}

	sql += ` order by sort desc, created_at desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
