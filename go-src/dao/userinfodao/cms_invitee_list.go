package userinfodao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	userentity "xr-game-server/entity/user"
)

// CMSInviteeListFilter CMS邀请下级用户查询条件
type CMSInviteeListFilter struct {
	InviterId uint64
	StartTime int64
	EndTime   int64
	PageIndex int
	PageSize  int
}

// CMSInviteeList 分页查询邀请下级(按资料创建时间倒序)
func CMSInviteeList(f *CMSInviteeListFilter) (int, []*userentity.UserInfo) {
	list := make([]*userentity.UserInfo, 0)
	if f == nil || f.InviterId == 0 {
		return 0, list
	}
	pageIndex := f.PageIndex
	if pageIndex <= 0 {
		pageIndex = 1
	}
	pageSize := f.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	ctx := gctx.New()
	m := g.Model(string(userentity.TbUserInfo)).Ctx(ctx).
		Where(string(userentity.UserInfoInviterId)+" = ?", f.InviterId)
	if f.StartTime > 0 {
		m = m.Where("created_at >= ?", time.Unix(f.StartTime, 0))
	}
	if f.EndTime > 0 {
		m = m.Where("created_at <= ?", time.Unix(f.EndTime, 0))
	}
	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	_ = m.Clone().Order("id desc").
		Limit(pageSize).Offset((pageIndex - 1) * pageSize).
		Scan(&list)
	return total, list
}
