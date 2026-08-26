package statdao

import (
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/entity/stat"
)

func PublishDailyLoginStat(data *entity.DailyLoginStat) {
	if data == nil || data.ID == "" || dailyLoginStatCacheMgr == nil {
		return
	}
	dailyLoginStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func PublishWeeklyLoginStat(data *entity.WeeklyLoginStat) {
	if data == nil || data.ID == "" || weeklyLoginStatCacheMgr == nil {
		return
	}
	weeklyLoginStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func PublishMonthlyLoginStat(data *entity.MonthlyLoginStat) {
	if data == nil || data.ID == "" || monthlyLoginStatCacheMgr == nil {
		return
	}
	monthlyLoginStatCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func publishDailyUserLogin(data *entity.DailyUserLogin) {
	if data == nil || data.ID == "" || dailyUserLoginCacheMgr == nil {
		return
	}
	dailyUserLoginCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func publishWeeklyUserLogin(data *entity.WeeklyUserLogin) {
	if data == nil || data.ID == "" || weeklyUserLoginCacheMgr == nil {
		return
	}
	weeklyUserLoginCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func publishMonthlyUserLogin(data *entity.MonthlyUserLogin) {
	if data == nil || data.ID == "" || monthlyUserLoginCacheMgr == nil {
		return
	}
	monthlyUserLoginCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func publishDailyUserAudience(data *entity.DailyUserAudience) {
	if data == nil || data.ID == "" || dailyUserAudienceCacheMgr == nil {
		return
	}
	dailyUserAudienceCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func publishWeeklyUserAudience(data *entity.WeeklyUserAudience) {
	if data == nil || data.ID == "" || weeklyUserAudienceCacheMgr == nil {
		return
	}
	weeklyUserAudienceCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func publishMonthlyUserAudience(data *entity.MonthlyUserAudience) {
	if data == nil || data.ID == "" || monthlyUserAudienceCacheMgr == nil {
		return
	}
	monthlyUserAudienceCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func publishDailyUserRecharge(data *entity.DailyUserRecharge) {
	if data == nil || data.ID == "" || dailyUserRechargeCacheMgr == nil {
		return
	}
	dailyUserRechargeCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func publishWeeklyUserRecharge(data *entity.WeeklyUserRecharge) {
	if data == nil || data.ID == "" || weeklyUserRechargeCacheMgr == nil {
		return
	}
	weeklyUserRechargeCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func publishMonthlyUserRecharge(data *entity.MonthlyUserRecharge) {
	if data == nil || data.ID == "" || monthlyUserRechargeCacheMgr == nil {
		return
	}
	monthlyUserRechargeCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func publishDailyUserGoldConsume(data *entity.DailyUserGoldConsume) {
	if data == nil || data.ID == "" || dailyUserGoldConsumeCacheMgr == nil {
		return
	}
	dailyUserGoldConsumeCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func publishWeeklyUserGoldConsume(data *entity.WeeklyUserGoldConsume) {
	if data == nil || data.ID == "" || weeklyUserGoldConsumeCacheMgr == nil {
		return
	}
	weeklyUserGoldConsumeCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func publishMonthlyUserGoldConsume(data *entity.MonthlyUserGoldConsume) {
	if data == nil || data.ID == "" || monthlyUserGoldConsumeCacheMgr == nil {
		return
	}
	monthlyUserGoldConsumeCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func publishDailyUserDiamondConsume(data *entity.DailyUserDiamondConsume) {
	if data == nil || data.ID == "" || dailyUserDiamondConsumeCacheMgr == nil {
		return
	}
	dailyUserDiamondConsumeCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func publishWeeklyUserDiamondConsume(data *entity.WeeklyUserDiamondConsume) {
	if data == nil || data.ID == "" || weeklyUserDiamondConsumeCacheMgr == nil {
		return
	}
	weeklyUserDiamondConsumeCacheMgr.PublishRow(gctx.New(), data.ID, data)
}

func publishMonthlyUserDiamondConsume(data *entity.MonthlyUserDiamondConsume) {
	if data == nil || data.ID == "" || monthlyUserDiamondConsumeCacheMgr == nil {
		return
	}
	monthlyUserDiamondConsumeCacheMgr.PublishRow(gctx.New(), data.ID, data)
}
