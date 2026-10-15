package userwithdrawdao

import (
	"sort"

	userentity "xr-game-server/entity/user"
)

func cachedUserWithdrawLogByID(userId, logId uint64) (*userentity.UserWithdrawLog, bool) {
	if userId == 0 || logId == 0 {
		return nil, false
	}
	list := ensureAppUserWithdrawLogListCacheLoaded(userId)
	for _, item := range list {
		if item != nil && item.ID == logId {
			return cloneUserWithdrawLog(item), true
		}
	}
	return nil, false
}

// PreferCachedUserWithdrawLog 先以库行为准; App 列表缓存有同 id 则用缓存行
func PreferCachedUserWithdrawLog(row *userentity.UserWithdrawLog) *userentity.UserWithdrawLog {
	if row == nil {
		return nil
	}
	if cached, ok := cachedUserWithdrawLogByID(row.UserId, row.ID); ok {
		return cached
	}
	return row
}

// MergeCMSUserWithdrawLogs CMS 列表查库后, 逐条 Prefer 缓存
func MergeCMSUserWithdrawLogs(filter *UserWithdrawListFilter, dbTotal int, dbRows []*userentity.UserWithdrawLog) (int, []*userentity.UserWithdrawLog) {
	if filter != nil && filter.UserId > 0 {
		ensureAppUserWithdrawLogListCacheLoaded(filter.UserId)
	}
	rows := make([]*userentity.UserWithdrawLog, 0, len(dbRows))
	for _, row := range dbRows {
		if row == nil {
			continue
		}
		rows = append(rows, PreferCachedUserWithdrawLog(row))
	}
	sortUserWithdrawLogsByUpdatedAtDesc(rows)
	return dbTotal, rows
}

func sortUserWithdrawLogsByUpdatedAtDesc(rows []*userentity.UserWithdrawLog) {
	sort.Slice(rows, func(i, j int) bool {
		a, b := rows[i], rows[j]
		if a == nil || b == nil {
			return a != nil
		}
		if a.UpdatedAt.Equal(b.UpdatedAt) {
			return a.ID > b.ID
		}
		return a.UpdatedAt.After(b.UpdatedAt)
	})
}
