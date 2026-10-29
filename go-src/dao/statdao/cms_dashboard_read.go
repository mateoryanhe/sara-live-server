package statdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/entity/stat"
)

// ScanSystemTotalStatForCMS 仪表盘组装时从库读累计统计(不替换进程内 sysSata, 避免与 syndb 写入打架).
func ScanSystemTotalStatForCMS() *entity.SystemTotalStat {
	var data *entity.SystemTotalStat
	_ = g.Model(string(entity.TbSystemTotalStat)).Limit(1).Scan(&data)
	if data != nil {
		return data
	}
	if live := GetSysStat(); live != nil {
		return live
	}
	return entity.NewSystemTotalStat(entity.SystemTotalStatDefaultID)
}

// ScanDailyLoginStatByDate CMS 仪表盘读今日聚合, 绕过 RowCache, 供 5 分钟整包 cache miss 时查库.
func ScanDailyLoginStatByDate(date string) *entity.DailyLoginStat {
	var data *entity.DailyLoginStat
	_ = g.Model(string(entity.TbDailyLoginStat)).Unscoped().Where(g.Map{
		string(db.IdName): date,
	}).Scan(&data)
	if data != nil {
		return data
	}
	return entity.NewDailyLoginStat(date)
}

// ScanWeeklyLoginStatByWeek CMS 仪表盘读本周聚合, 绕过 RowCache.
func ScanWeeklyLoginStatByWeek(week string) *entity.WeeklyLoginStat {
	var data *entity.WeeklyLoginStat
	_ = g.Model(string(entity.TbWeeklyLoginStat)).Unscoped().Where(g.Map{
		string(db.IdName): week,
	}).Scan(&data)
	if data != nil {
		return data
	}
	return entity.NewWeeklyLoginStat(week)
}

// ScanMonthlyLoginStatByMonth CMS 仪表盘读本月聚合, 绕过 RowCache.
func ScanMonthlyLoginStatByMonth(month string) *entity.MonthlyLoginStat {
	var data *entity.MonthlyLoginStat
	_ = g.Model(string(entity.TbMonthlyLoginStat)).Unscoped().Where(g.Map{
		string(db.IdName): month,
	}).Scan(&data)
	if data != nil {
		return data
	}
	return entity.NewMonthlyLoginStat(month)
}
