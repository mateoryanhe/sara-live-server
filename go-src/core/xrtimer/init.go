package xrtimer

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/util/gutil"
	"time"
)

func Init() {

}

// AddSingleton 执行定时任务,可以取消,有防panic,避免程序崩溃,执行完任务,才会执行下个任务,长时间间隔的定时任务,使用此方法
// 玩家的技能升级时间,cd时间使用,定时1个小时,每周凌晨0点执行,都可以使用
// job: 函数
// interval: 间隔时间 大于秒的定时任务可以使用此方法
func AddSingleton(ctx context.Context, interval time.Duration, job gtimer.JobFunc) *gtimer.Entry {
	entry := gtimer.AddSingleton(ctx, interval, func(ctx context.Context) {
		gutil.TryCatch(ctx, func(try context.Context) {
			job(ctx)
		}, func(catch context.Context, exception error) {
			g.Log().Errorf(catch, "循环-定时器报错 AddSingleton error %v", gerror.Stack(exception))
		})
	})
	return entry
}

// AddOnce 只执行一次,可以取消,有防panic,避免程序崩溃,时间周期小于秒,使用,秒级任务请使用AddLoopTask
// 隔间时间长的定时任务可以使用此方法,可以高效利用cpu,避免cpu一直在循环,占用cpu,
// 玩家的技能升级时间,cd时间使用,定时1个小时,每周凌晨0点执行,都可以使用
// job: 任务
// interval: 任务时间
func AddOnce(ctx context.Context, taskTime time.Duration, job gtimer.JobFunc) *gtimer.Entry {

	entry := gtimer.AddOnce(ctx, taskTime, func(ctx context.Context) {
		gutil.TryCatch(ctx, func(try context.Context) {
			job(ctx)
		}, func(catch context.Context, exception error) {
			g.Log().Errorf(catch, "单次-定时器报错 AddOnce error %v", gerror.Stack(exception))
		})
	})
	return entry
}
