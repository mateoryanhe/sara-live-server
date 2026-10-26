package userinfodao

import (
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	userentity "xr-game-server/entity/user"
)

// CMSInviteLogListFilter CMS全局邀请日志查询条件
type CMSInviteLogListFilter struct {
	Keyword   string
	StartTime int64
	EndTime   int64
	PageIndex int
	PageSize  int
}

// CMSInviteLogList 分页查询有邀请人的用户(按资料创建时间倒序)
func CMSInviteLogList(f *CMSInviteLogListFilter) (int, []*userentity.UserInfo) {
	list := make([]*userentity.UserInfo, 0)
	if f == nil {
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
	keyword := strings.TrimSpace(f.Keyword)
	m := g.Model(string(userentity.TbUserInfo)+" u").Ctx(ctx).
		Where("u."+string(userentity.UserInfoInviterId)+" > ?", 0)
	if keyword != "" {
		like := "%" + keyword + "%"
		m = m.LeftJoin(string(userentity.TbUserInfo)+" inv", "inv.id = u."+string(userentity.UserInfoInviterId)).
			Where("(CAST(u.id AS CHAR) LIKE ? OR u."+string(userentity.UserInfoNickname)+" LIKE ? OR CAST(u."+string(userentity.UserInfoInviterId)+" AS CHAR) LIKE ? OR inv."+string(userentity.UserInfoNickname)+" LIKE ?)",
				like, like, like, like)
	}
	if f.StartTime > 0 {
		m = m.Where("u.created_at >= ?", time.Unix(f.StartTime, 0))
	}
	if f.EndTime > 0 {
		m = m.Where("u.created_at <= ?", time.Unix(f.EndTime, 0))
	}
	total, err := m.Clone().Count()
	if err != nil {
		return 0, list
	}
	_ = m.Clone().Fields("u.*").Order("u.id desc").
		Limit(pageSize).Offset((pageIndex - 1) * pageSize).
		Scan(&list)
	return total, list
}
