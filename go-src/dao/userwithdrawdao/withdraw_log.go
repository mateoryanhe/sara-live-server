package userwithdrawdao

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	userentity "xr-game-server/entity/user"
)

type UserWithdrawListFilter struct {
	UserId   uint64
	Status   *uint8
	Page     int
	PageSize int
}

func GetUserWithdrawLogById(id uint64) *userentity.UserWithdrawLog {
	if id == 0 {
		return nil
	}
	var row userentity.UserWithdrawLog
	err := g.DB().Model(string(userentity.TbUserWithdrawLog)).WherePri(id).Scan(&row)
	if err != nil || row.ID == 0 {
		return nil
	}
	return PreferCachedUserWithdrawLog(&row)
}

// PublishNewUserWithdrawLog 通过 entity setter 已入 syndb 队列, 并刷新 App 列表缓存
func PublishNewUserWithdrawLog(row *userentity.UserWithdrawLog) {
	PrependUserWithdrawToAppListCache(row)
}

// SyncUserWithdrawLogCache 修改后 syndb 已排队; 按 userId 刷新列表缓存并将该单置顶
func SyncUserWithdrawLogCache(row *userentity.UserWithdrawLog) {
	PrependUserWithdrawToAppListCache(row)
}

func ListUserWithdrawLogs(filter *UserWithdrawListFilter) (int, []*userentity.UserWithdrawLog, error) {
	rows := make([]*userentity.UserWithdrawLog, 0)
	if filter == nil {
		return 0, rows, nil
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	model := g.DB().Model(string(userentity.TbUserWithdrawLog)).Ctx(gctx.New())
	if filter.UserId > 0 {
		model = model.Where(string(userentity.UserWithdrawLogUserId)+" = ?", filter.UserId)
	}
	if filter.Status != nil {
		model = model.Where(string(userentity.UserWithdrawLogStatus)+" = ?", *filter.Status)
	}
	total, err := model.Count()
	if err != nil {
		return 0, nil, err
	}
	err = model.Order(userWithdrawLogListOrder).
		Page(page, pageSize).
		Scan(&rows)
	if err != nil {
		return 0, nil, err
	}
	total, rows = MergeCMSUserWithdrawLogs(filter, total, rows)
	return total, rows, nil
}

func CountUserWithdrawLogsSince(userId uint64, since time.Time) (int, error) {
	if userId == 0 {
		return 0, nil
	}
	n, err := g.DB().Model(string(userentity.TbUserWithdrawLog)).Ctx(gctx.New()).
		Where(string(userentity.UserWithdrawLogUserId)+" = ?", userId).
		Where("created_at >= ?", since).
		Count()
	if err != nil {
		return 0, err
	}
	return n, nil
}
