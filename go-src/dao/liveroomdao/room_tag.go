package liveroomdao

import (
	"strconv"
	"xr-game-server/core/str"
	"xr-game-server/dto/liveroomdto"
	"xr-game-server/entity"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func GetRoomTagById(id uint64) *entity.LiveRoomTag {
	var row entity.LiveRoomTag
	if err := g.DB().Model(string(entity.TbLiveRoomTag)).Where("id = ?", id).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func GetRoomTagByName(name string) *entity.LiveRoomTag {
	if name == "" {
		return nil
	}
	var row entity.LiveRoomTag
	if err := g.DB().Model(string(entity.TbLiveRoomTag)).Where("name = ?", name).Scan(&row); err != nil {
		return nil
	}
	if row.ID == 0 {
		return nil
	}
	return &row
}

func GetAllRoomTags() []*entity.LiveRoomTag {
	ret := make([]*entity.LiveRoomTag, 0)
	_ = g.DB().Model(string(entity.TbLiveRoomTag)).
		Order("sort desc, created_at desc").
		Scan(&ret)
	return ret
}

func CreateRoomTag(row *entity.LiveRoomTag) error {
	_, err := g.DB().Model(string(entity.TbLiveRoomTag)).Save(row)
	return err
}

func UpdateRoomTag(row *entity.LiveRoomTag) error {
	return CreateRoomTag(row)
}

func DeleteRoomTag(id uint64) error {
	_, err := g.DB().Model(string(entity.TbLiveRoomTag)).WherePri(id).Delete()
	return err
}

func GetRoomTagList(req *liveroomdto.LiveRoomTagListReq) (int, []*liveroomdto.LiveRoomTagListRes) {
	sql := `select id, name, sort, created_at, updated_at
            from live_room_tags
            where 1=1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*liveroomdto.LiveRoomTagListRes, 0)

	sql += ` order by sort desc, created_at desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` limit ` + strconv.Itoa(req.PageSize) + ` offset ` + strconv.Itoa(req.PageIndex-1)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
