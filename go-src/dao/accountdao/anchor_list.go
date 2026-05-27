package accountdao

import (
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/str"
	"xr-game-server/dto/accountdto"
)

// GetAnchorList CMS分页查询主播列表(基于 user_infos.is_anchor = 1)
func GetAnchorList(req *accountdto.QueryAnchorListReq) (int, []*accountdto.AnchorListItem) {
	sql := `select u.id,
                   u.nickname, u.phone, u.avatar, u.guild_id,
                   u.created_at,
                   a.ip, a.created_at as registered_at,
                   ifnull(r.ban, 0) as ban, r.ban_apply_time, ifnull(r.ban_reason, '') as ban_reason,
                   ifnull(r.title, '') as room_title,
                   if(r.live_record_id > 0, 1, 0) as live_status
            from user_infos u
            inner join accounts a on a.id = u.id
            left join live_rooms r on r.id = u.id
            where u.is_anchor = 1 `
	param := make([]any, 0)
	ctx := gctx.New()
	ret := make([]*accountdto.AnchorListItem, 0)

	if req.Key != "" {
		sql += ` and (u.id = ? or u.nickname like ? or u.phone like ? or u.share_code like ?)`
		likeKey := "%" + req.Key + "%"
		param = append(param, req.Key, likeKey, likeKey, likeKey)
	}

	sql += ` order by u.id desc`
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)

	pageIndex := req.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (pageIndex - 1) * pageSize
	sql += ` limit ` + strconv.Itoa(pageSize) + ` offset ` + strconv.Itoa(offset)
	_ = g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
