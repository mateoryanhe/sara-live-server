package statdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/entity/stat"
)

func ListDailyLoginStatsByIDs(ids []string) []*entity.DailyLoginStat {
	if len(ids) == 0 {
		return nil
	}
	rows := make([]*entity.DailyLoginStat, 0, len(ids))
	_ = g.Model(string(entity.TbDailyLoginStat)).WhereIn(string(db.IdName), ids).Scan(&rows)
	return orderDailyLoginStats(ids, rows)
}

func ListWeeklyLoginStatsByIDs(ids []string) []*entity.WeeklyLoginStat {
	if len(ids) == 0 {
		return nil
	}
	rows := make([]*entity.WeeklyLoginStat, 0, len(ids))
	_ = g.Model(string(entity.TbWeeklyLoginStat)).WhereIn(string(db.IdName), ids).Scan(&rows)
	return orderWeeklyLoginStats(ids, rows)
}

func ListMonthlyLoginStatsByIDs(ids []string) []*entity.MonthlyLoginStat {
	if len(ids) == 0 {
		return nil
	}
	rows := make([]*entity.MonthlyLoginStat, 0, len(ids))
	_ = g.Model(string(entity.TbMonthlyLoginStat)).WhereIn(string(db.IdName), ids).Scan(&rows)
	return orderMonthlyLoginStats(ids, rows)
}

func orderDailyLoginStats(ids []string, rows []*entity.DailyLoginStat) []*entity.DailyLoginStat {
	byID := make(map[string]*entity.DailyLoginStat, len(rows))
	for _, row := range rows {
		if row != nil && row.ID != "" {
			byID[row.ID] = row
		}
	}
	out := make([]*entity.DailyLoginStat, 0, len(ids))
	for _, id := range ids {
		if row, ok := byID[id]; ok {
			out = append(out, row)
			continue
		}
		out = append(out, entity.NewDailyLoginStat(id))
	}
	return out
}

func orderWeeklyLoginStats(ids []string, rows []*entity.WeeklyLoginStat) []*entity.WeeklyLoginStat {
	byID := make(map[string]*entity.WeeklyLoginStat, len(rows))
	for _, row := range rows {
		if row != nil && row.ID != "" {
			byID[row.ID] = row
		}
	}
	out := make([]*entity.WeeklyLoginStat, 0, len(ids))
	for _, id := range ids {
		if row, ok := byID[id]; ok {
			out = append(out, row)
			continue
		}
		out = append(out, entity.NewWeeklyLoginStat(id))
	}
	return out
}

func orderMonthlyLoginStats(ids []string, rows []*entity.MonthlyLoginStat) []*entity.MonthlyLoginStat {
	byID := make(map[string]*entity.MonthlyLoginStat, len(rows))
	for _, row := range rows {
		if row != nil && row.ID != "" {
			byID[row.ID] = row
		}
	}
	out := make([]*entity.MonthlyLoginStat, 0, len(ids))
	for _, id := range ids {
		if row, ok := byID[id]; ok {
			out = append(out, row)
			continue
		}
		out = append(out, entity.NewMonthlyLoginStat(id))
	}
	return out
}
