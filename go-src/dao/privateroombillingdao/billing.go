package privateroombillingdao

import (
	"strconv"
	"xr-game-server/core/str"
	"xr-game-server/dto/privateroombillingdto"
	"xr-game-server/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func GetById(id uint64) *entity.LivePrivateRoomBilling {
	var row entity.LivePrivateRoomBilling
	if err := g.DB().Model(string(entity.TbLivePrivateRoomBilling)).Where("id = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func GetAll() []*entity.LivePrivateRoomBilling {
	ret := make([]*entity.LivePrivateRoomBilling, 0)
	_ = g.DB().Model(string(entity.TbLivePrivateRoomBilling)).
		Order("sort desc, created_at desc").
		Scan(&ret)
	return ret
}

func Create(row *entity.LivePrivateRoomBilling) error {
	_, err := g.DB().Model(string(entity.TbLivePrivateRoomBilling)).Save(row)
	return err
}

func Update(row *entity.LivePrivateRoomBilling) error {
	return Create(row)
}

func Delete(id uint64) error {
	_, err := g.DB().Model(string(entity.TbLivePrivateRoomBilling)).WherePri(id).Delete()
	return err
}

func UpdateStatus(id uint64, status uint8) error {
	_, err := g.DB().Model(string(entity.TbLivePrivateRoomBilling)).
		WherePri(id).
		Data(g.Map{"status": status}).
		Update()
	return err
}

func GetBillingList(req *privateroombillingdto.BillingListReq) (int, []*privateroombillingdto.BillingListRes) {
	sql := `select id, price_per_minute, sort, status, created_at, updated_at
            from live_private_room_billings
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*privateroombillingdto.BillingListRes, 0)

	switch req.StatusFilter {
	case 1:
		sql += ` and status = ?`
		param = append(param, entity.LivePrivateRoomBillingStatusOffShelf)
	case 2:
		sql += ` and status = ?`
		param = append(param, entity.LivePrivateRoomBillingStatusOnShelf)
	}

	sql += ` order by sort desc, created_at desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
