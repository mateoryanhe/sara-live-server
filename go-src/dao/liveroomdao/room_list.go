package liveroomdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"strconv"
	"xr-game-server/core/str"
)

// RoomListRow 直播间列表查询行(直连 DB)
type RoomListRow struct {
	ID        uint64 `json:"id"`
	GuildId   uint64 `json:"guildId"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Notice    string `json:"notice"`
	Status    uint8  `json:"status"`
	CreatedAt int64  `json:"createdAt"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
}

// ListRooms 分页查询直播间列表(直连数据库,不走缓存)
// statusFilter: 0=全部, 1=仅直播中, 2=仅未开播/已下播
func ListRooms(page, pageSize int) (int, []*RoomListRow) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	sql := `SELECT r.id, r.guild_id, r.title, r.cover, r.notice, r.status,
                   UNIX_TIMESTAMP(r.created_at) AS created_at,
                   IFNULL(u.nickname, '') AS nickname,
                   IFNULL(u.avatar, '') AS avatar
            FROM live_rooms r
            LEFT JOIN user_infos u ON u.id = r.id
            WHERE 1=1 `
	param := make([]any, 0)

	sql += ` ORDER BY r.status DESC, r.updated_at DESC`

	ctx := gctx.New()
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` LIMIT ` + strconv.Itoa(pageSize) + ` OFFSET ` + strconv.Itoa((page-1)*pageSize)

	ret := make([]*RoomListRow, 0)
	g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
