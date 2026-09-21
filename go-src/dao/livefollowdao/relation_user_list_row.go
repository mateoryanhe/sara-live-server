package livefollowdao

import (
	"time"

	"xr-game-server/entity/live"
	userentity "xr-game-server/entity/user"
)

// RelationUserListRow 是关注关系与列表目标用户资料的聚合缓存行。
// 粉丝列表中的用户资料对应 UserId；黑名单列表中的用户资料对应 AnchorId。
type RelationUserListRow struct {
	ID        string     `orm:"id"`
	UserId    uint64     `orm:"user_id"`
	AnchorId  uint64     `orm:"anchor_id"`
	Status    uint8      `orm:"status"`
	CreatedAt time.Time  `orm:"created_at"`
	UpdatedAt time.Time  `orm:"updated_at"`
	Nickname  string     `orm:"nickname"`
	Avatar    string     `orm:"avatar"`
	VipLevel  uint32     `orm:"vip_level"`
	Gender    uint8      `orm:"gender"`
	Birthday  *time.Time `orm:"birthday"`
}

func newRelationUserListRow(f *entity.LiveFollow, user *userentity.UserInfo) *RelationUserListRow {
	if f == nil {
		return nil
	}
	row := &RelationUserListRow{
		ID:        f.ID,
		UserId:    f.UserId,
		AnchorId:  f.AnchorId,
		Status:    f.Status,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
	}
	if user != nil {
		row.Nickname = user.Nickname
		row.Avatar = user.Avatar
		row.VipLevel = user.VipLevel
		row.Gender = user.Gender
		row.Birthday = user.Birthday
	}
	return row
}

func (row *RelationUserListRow) liveFollow() *entity.LiveFollow {
	if row == nil {
		return nil
	}
	return &entity.LiveFollow{
		ID:        row.ID,
		UserId:    row.UserId,
		AnchorId:  row.AnchorId,
		Status:    row.Status,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func relationUserListFields(relationAlias, userAlias string) []any {
	return []any{
		relationAlias + ".id AS id",
		relationAlias + ".user_id AS user_id",
		relationAlias + ".anchor_id AS anchor_id",
		relationAlias + ".status AS status",
		relationAlias + ".created_at AS created_at",
		relationAlias + ".updated_at AS updated_at",
		userAlias + ".nickname AS nickname",
		userAlias + ".avatar AS avatar",
		userAlias + ".vip_level AS vip_level",
		userAlias + ".gender AS gender",
		userAlias + ".birthday AS birthday",
	}
}
