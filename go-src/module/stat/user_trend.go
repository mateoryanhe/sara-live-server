package stat

import (
	"context"
	"xr-game-server/dao/dailyloginstatdao"
	"xr-game-server/dao/monthlyloginstatdao"
	"xr-game-server/dao/weeklyloginstatdao"
	"xr-game-server/dto/statdto"
	"xr-game-server/entity"
)

const (
	userStatDailyLimit   = 30
	userStatWeeklyLimit  = 12
	userStatMonthlyLimit = 12
)

// GetCMSUserStatTrend CMS获取用户数据趋势
func GetCMSUserStatTrend(_ context.Context, _ *statdto.CMSUserStatTrendReq) (*statdto.CMSUserStatTrendRes, error) {
	return &statdto.CMSUserStatTrendRes{
		Daily:   toDailyTrendPoints(dailyloginstatdao.ListRecent(userStatDailyLimit)),
		Weekly:  toWeeklyTrendPoints(weeklyloginstatdao.ListRecent(userStatWeeklyLimit)),
		Monthly: toMonthlyTrendPoints(monthlyloginstatdao.ListRecent(userStatMonthlyLimit)),
	}, nil
}

func toDailyTrendPoints(rows []*entity.DailyLoginStat) []*statdto.CMSUserStatTrendPoint {
	list := make([]*statdto.CMSUserStatTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &statdto.CMSUserStatTrendPoint{
			Time:              row.ID,
			ActiveUserCount:   row.Count,
			RegisterUserCount: row.RegisterCount,
		})
	}
	return list
}

func toWeeklyTrendPoints(rows []*entity.WeeklyLoginStat) []*statdto.CMSUserStatTrendPoint {
	list := make([]*statdto.CMSUserStatTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &statdto.CMSUserStatTrendPoint{
			Time:              row.ID,
			ActiveUserCount:   row.Count,
			RegisterUserCount: row.RegisterCount,
		})
	}
	return list
}

func toMonthlyTrendPoints(rows []*entity.MonthlyLoginStat) []*statdto.CMSUserStatTrendPoint {
	list := make([]*statdto.CMSUserStatTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		list = append(list, &statdto.CMSUserStatTrendPoint{
			Time:              row.ID,
			ActiveUserCount:   row.Count,
			RegisterUserCount: row.RegisterCount,
		})
	}
	return list
}
