package entryeffectdao

import (
	"fmt"
	"strconv"
	"xr-game-server/core/str"
	"xr-game-server/dto/entryeffectdto"
	"xr-game-server/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func GetById(id uint64) *entity.LiveEntryEffect {
	var row entity.LiveEntryEffect
	if err := g.DB().Model(string(entity.TbLiveEntryEffect)).Where("id = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func GetByName(name string) *entity.LiveEntryEffect {
	var row entity.LiveEntryEffect
	if err := g.DB().Model(string(entity.TbLiveEntryEffect)).Where("name = ?", name).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func GetAll() []*entity.LiveEntryEffect {
	ret := make([]*entity.LiveEntryEffect, 0)
	_ = g.DB().Model(string(entity.TbLiveEntryEffect)).
		Order("level_start asc, level_end asc, created_at desc").
		Scan(&ret)
	return ret
}

func Create(row *entity.LiveEntryEffect) error {
	_, err := g.DB().Model(string(entity.TbLiveEntryEffect)).Save(row)
	return err
}

func Update(row *entity.LiveEntryEffect) error {
	return Create(row)
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbLiveEntryEffect)).WherePri(id).Delete()
	return err
}

func UpdateStatus(id uint64, status uint8) error {
	_, err := g.DB().Model(string(entity.TbLiveEntryEffect)).
		WherePri(id).
		Data(g.Map{"status": status}).
		Update()
	return err
}

func GetEntryEffectList(req *entryeffectdto.EntryEffectListReq) (int, []*entryeffectdto.EntryEffectListRes) {
	sql := `select id, name, level_start, level_end, animation, status, created_at, updated_at
            from live_entry_effects
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*entryeffectdto.EntryEffectListRes, 0)

	if req.Name != "" {
		sql += ` and name LIKE ?`
		param = append(param, fmt.Sprintf("%%%s%%", req.Name))
	}
	switch req.StatusFilter {
	case 1:
		sql += ` and status = ?`
		param = append(param, entity.LiveEntryEffectStatusOffShelf)
	case 2:
		sql += ` and status = ?`
		param = append(param, entity.LiveEntryEffectStatusOnShelf)
	}

	sql += ` order by level_start asc, level_end asc, created_at desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
