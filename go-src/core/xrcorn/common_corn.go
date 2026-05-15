package xrcorn

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

const (
	DayPattern   = "0 0 0 * * *"
	WeekPattern  = "0 0 0 * * 1"
	MonthPattern = "0 0 0 1 * *"
)

func initDay() {
	// 添加定时任务，每天0点执行
	next, _ := GetNextTime(DayPattern)
	g.Log().Infof(gctx.New(), "day执行下次任务时间:%v", next)
	_, _ = AddSingleton(gctx.New(), DayPattern, func(ctx context.Context) {
		g.Log().Infof(ctx, "开始执行每日0点任务")
		event.Pub(gameevent.DayEvent, gameevent.NewDayEventData())
		g.Log().Infof(ctx, "执行每日0点任务完毕")
		newNext, _ := GetNextTime(DayPattern)
		g.Log().Infof(gctx.New(), "day执行下次任务时间:%v", newNext)
	})
}

/*
0 0 0 * * 1 各字段含义如下：

第 1 位（秒）：0 → 0 秒
第 2 位（分）：0 → 0 分
第 3 位（时）：0 → 0 点
第 4 位（日）：* → 任意日期
第 5 位（月）：* → 任意月份
第 6 位（周）：1 → 周一（取值范围 0-6，0 和 7 都代表周日）
*/
func initWeek() {
	// 添加定时任务，每周一0点执行
	next, _ := GetNextTime(WeekPattern)
	g.Log().Infof(gctx.New(), "week执行下次任务时间:%v", next)
	_, _ = AddSingleton(gctx.New(), WeekPattern, func(ctx context.Context) {
		g.Log().Infof(ctx, "开始执行每周一0点任务")
		event.Pub(gameevent.WeekEvent, gameevent.NewWeekEventData())
		g.Log().Infof(ctx, "执行每周一0点任务完毕")
		newNext, _ := GetNextTime(WeekPattern)
		g.Log().Infof(gctx.New(), "week执行下次任务时间:%v", newNext)
	})
}

/*
0 0 0 1 * * 各字段含义如下：

第 1 位（秒）：0 → 0 秒
第 2 位（分）：0 → 0 分
第 3 位（时）：0 → 0 点
第 4 位（日）：1 → 每月的第一天
第 5 位（月）：* → 任意月份
第 6 位（周）：* → 任意星期
*/
func initMonth() {
	// 添加定时任务，每月第一天0点执行
	next, _ := GetNextTime(MonthPattern)
	g.Log().Infof(gctx.New(), "month执行下次任务时间:%v", next)
	_, _ = AddSingleton(gctx.New(), MonthPattern, func(ctx context.Context) {
		g.Log().Infof(ctx, "开始执行每月第一天0点任务")
		event.Pub(gameevent.MonthEvent, gameevent.NewMonthEventData())
		g.Log().Infof(ctx, "执行每月第一天0点任务完毕")
		nextMonth, _ := GetNextTime(MonthPattern)
		g.Log().Infof(gctx.New(), "month执行下次任务时间:%v", nextMonth)
	})

}
