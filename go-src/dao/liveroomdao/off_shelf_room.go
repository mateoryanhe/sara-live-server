package liveroomdao

import (
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/str"
	"xr-game-server/entity"
)

// OffShelfRoomRow 回收站列表行(直连 DB)
type OffShelfRoomRow struct {
	ID           uint64     `json:"id"`
	GuildId      uint64     `json:"guildId"`
	Title        string     `json:"title"`
	Cover        string     `json:"cover"`
	Category     uint8      `json:"category"`
	Nickname     string     `json:"nickname"`
	Avatar       string     `json:"avatar"`
	Phone        string     `json:"phone"`
	UserType     uint8      `json:"userType"`
	Ban          bool       `json:"ban"`
	BanApplyTime *time.Time `json:"banApplyTime"`
	BanReason    string     `json:"banReason"`
	UpdatedAt    *time.Time `json:"updatedAt"`
	CreatedAt    *time.Time `json:"createdAt"`
}

// ListOffShelfRooms 分页查询已下架直播间(直连数据库,不走缓存)
func ListOffShelfRooms(page, pageSize int, key string) (int, []*OffShelfRoomRow) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	sql := `SELECT r.id, r.guild_id, r.title, r.cover, r.category,
                   r.ban, r.ban_apply_time, r.ban_reason, r.updated_at, r.created_at,
                   IFNULL(u.nickname, '') AS nickname,
                   IFNULL(u.avatar, '') AS avatar,
                   IFNULL(u.phone, '') AS phone,
                   IFNULL(u.user_type, 0) AS user_type
            FROM live_rooms r
            LEFT JOIN user_infos u ON u.id = r.id
            WHERE r.status = ?`
	param := []any{entity.LiveRoomStatusOffShelf}

	key = strings.TrimSpace(key)
	if key != "" {
		like := "%" + key + "%"
		sql += ` AND (CAST(r.id AS CHAR) LIKE ? OR u.nickname LIKE ? OR u.phone LIKE ? OR u.share_code LIKE ?)`
		param = append(param, like, like, like, like)
	}

	sql += ` ORDER BY r.updated_at DESC`

	ctx := gctx.New()
	countSql := str.GetCountSQL(sql)
	total, _ := g.DB().GetCount(ctx, countSql, param)
	sql += ` LIMIT ` + strconv.Itoa(pageSize) + ` OFFSET ` + strconv.Itoa((page-1)*pageSize)

	ret := make([]*OffShelfRoomRow, 0)
	_ = g.DB().GetScan(ctx, &ret, sql, param)
	return total, ret
}
