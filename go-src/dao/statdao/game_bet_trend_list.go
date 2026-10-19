package statdao

import (
	"github.com/gogf/gf/v2/frame/g"
	"xr-game-server/constants/db"
	"xr-game-server/entity/stat"
)

func ListDailyGameBetUserStatsByIDs(ids []string) []*entity.DailyGameBetUserStat {
	return listDailyGameBetUserStatsByIDs(ids)
}

func ListDailyGameOnlinePeakStatsByIDs(ids []string) []*entity.DailyGameOnlinePeakStat {
	return listDailyGameOnlinePeakStatsByIDs(ids)
}

func ListWeeklyGameBetUserStatsByIDs(ids []string) []*entity.WeeklyGameBetUserStat {
	return listWeeklyGameBetUserStatsByIDs(ids)
}

func ListWeeklyGameOnlinePeakStatsByIDs(ids []string) []*entity.WeeklyGameOnlinePeakStat {
	return listWeeklyGameOnlinePeakStatsByIDs(ids)
}

func ListMonthlyGameBetUserStatsByIDs(ids []string) []*entity.MonthlyGameBetUserStat {
	return listMonthlyGameBetUserStatsByIDs(ids)
}

func ListMonthlyGameOnlinePeakStatsByIDs(ids []string) []*entity.MonthlyGameOnlinePeakStat {
	return listMonthlyGameOnlinePeakStatsByIDs(ids)
}

func listDailyGameBetUserStatsByIDs(ids []string) []*entity.DailyGameBetUserStat {
	if len(ids) == 0 {
		return nil
	}
	rows := make([]*entity.DailyGameBetUserStat, 0, len(ids))
	_ = g.Model(string(entity.TbDailyGameBetUserStat)).WhereIn(string(db.IdName), ids).Scan(&rows)
	return orderDailyGameBetUserStats(ids, rows)
}

func listDailyGameOnlinePeakStatsByIDs(ids []string) []*entity.DailyGameOnlinePeakStat {
	if len(ids) == 0 {
		return nil
	}
	rows := make([]*entity.DailyGameOnlinePeakStat, 0, len(ids))
	_ = g.Model(string(entity.TbDailyGameOnlinePeakStat)).WhereIn(string(db.IdName), ids).Scan(&rows)
	return orderDailyGameOnlinePeakStats(ids, rows)
}

func listWeeklyGameBetUserStatsByIDs(ids []string) []*entity.WeeklyGameBetUserStat {
	if len(ids) == 0 {
		return nil
	}
	rows := make([]*entity.WeeklyGameBetUserStat, 0, len(ids))
	_ = g.Model(string(entity.TbWeeklyGameBetUserStat)).WhereIn(string(db.IdName), ids).Scan(&rows)
	return orderWeeklyGameBetUserStats(ids, rows)
}

func listWeeklyGameOnlinePeakStatsByIDs(ids []string) []*entity.WeeklyGameOnlinePeakStat {
	if len(ids) == 0 {
		return nil
	}
	rows := make([]*entity.WeeklyGameOnlinePeakStat, 0, len(ids))
	_ = g.Model(string(entity.TbWeeklyGameOnlinePeakStat)).WhereIn(string(db.IdName), ids).Scan(&rows)
	return orderWeeklyGameOnlinePeakStats(ids, rows)
}

func listMonthlyGameBetUserStatsByIDs(ids []string) []*entity.MonthlyGameBetUserStat {
	if len(ids) == 0 {
		return nil
	}
	rows := make([]*entity.MonthlyGameBetUserStat, 0, len(ids))
	_ = g.Model(string(entity.TbMonthlyGameBetUserStat)).WhereIn(string(db.IdName), ids).Scan(&rows)
	return orderMonthlyGameBetUserStats(ids, rows)
}

func listMonthlyGameOnlinePeakStatsByIDs(ids []string) []*entity.MonthlyGameOnlinePeakStat {
	if len(ids) == 0 {
		return nil
	}
	rows := make([]*entity.MonthlyGameOnlinePeakStat, 0, len(ids))
	_ = g.Model(string(entity.TbMonthlyGameOnlinePeakStat)).WhereIn(string(db.IdName), ids).Scan(&rows)
	return orderMonthlyGameOnlinePeakStats(ids, rows)
}

func orderDailyGameBetUserStats(ids []string, rows []*entity.DailyGameBetUserStat) []*entity.DailyGameBetUserStat {
	byID := make(map[string]*entity.DailyGameBetUserStat, len(rows))
	for _, row := range rows {
		if row != nil && row.ID != "" {
			byID[row.ID] = row
		}
	}
	out := make([]*entity.DailyGameBetUserStat, 0, len(ids))
	for _, id := range ids {
		if row, ok := byID[id]; ok {
			out = append(out, row)
			continue
		}
		out = append(out, entity.NewDailyGameBetUserStat(id))
	}
	return out
}

func orderDailyGameOnlinePeakStats(ids []string, rows []*entity.DailyGameOnlinePeakStat) []*entity.DailyGameOnlinePeakStat {
	byID := make(map[string]*entity.DailyGameOnlinePeakStat, len(rows))
	for _, row := range rows {
		if row != nil && row.ID != "" {
			byID[row.ID] = row
		}
	}
	out := make([]*entity.DailyGameOnlinePeakStat, 0, len(ids))
	for _, id := range ids {
		if row, ok := byID[id]; ok {
			out = append(out, row)
			continue
		}
		out = append(out, entity.NewDailyGameOnlinePeakStat(id))
	}
	return out
}

func orderWeeklyGameBetUserStats(ids []string, rows []*entity.WeeklyGameBetUserStat) []*entity.WeeklyGameBetUserStat {
	byID := make(map[string]*entity.WeeklyGameBetUserStat, len(rows))
	for _, row := range rows {
		if row != nil && row.ID != "" {
			byID[row.ID] = row
		}
	}
	out := make([]*entity.WeeklyGameBetUserStat, 0, len(ids))
	for _, id := range ids {
		if row, ok := byID[id]; ok {
			out = append(out, row)
			continue
		}
		out = append(out, entity.NewWeeklyGameBetUserStat(id))
	}
	return out
}

func orderWeeklyGameOnlinePeakStats(ids []string, rows []*entity.WeeklyGameOnlinePeakStat) []*entity.WeeklyGameOnlinePeakStat {
	byID := make(map[string]*entity.WeeklyGameOnlinePeakStat, len(rows))
	for _, row := range rows {
		if row != nil && row.ID != "" {
			byID[row.ID] = row
		}
	}
	out := make([]*entity.WeeklyGameOnlinePeakStat, 0, len(ids))
	for _, id := range ids {
		if row, ok := byID[id]; ok {
			out = append(out, row)
			continue
		}
		out = append(out, entity.NewWeeklyGameOnlinePeakStat(id))
	}
	return out
}

func orderMonthlyGameBetUserStats(ids []string, rows []*entity.MonthlyGameBetUserStat) []*entity.MonthlyGameBetUserStat {
	byID := make(map[string]*entity.MonthlyGameBetUserStat, len(rows))
	for _, row := range rows {
		if row != nil && row.ID != "" {
			byID[row.ID] = row
		}
	}
	out := make([]*entity.MonthlyGameBetUserStat, 0, len(ids))
	for _, id := range ids {
		if row, ok := byID[id]; ok {
			out = append(out, row)
			continue
		}
		out = append(out, entity.NewMonthlyGameBetUserStat(id))
	}
	return out
}

func orderMonthlyGameOnlinePeakStats(ids []string, rows []*entity.MonthlyGameOnlinePeakStat) []*entity.MonthlyGameOnlinePeakStat {
	byID := make(map[string]*entity.MonthlyGameOnlinePeakStat, len(rows))
	for _, row := range rows {
		if row != nil && row.ID != "" {
			byID[row.ID] = row
		}
	}
	out := make([]*entity.MonthlyGameOnlinePeakStat, 0, len(ids))
	for _, id := range ids {
		if row, ok := byID[id]; ok {
			out = append(out, row)
			continue
		}
		out = append(out, entity.NewMonthlyGameOnlinePeakStat(id))
	}
	return out
}
