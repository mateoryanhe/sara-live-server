package xrcorn

import (
	"github.com/robfig/cron/v3"
	"time"
)

// GetNextTime 根据cron表达式获取下次执行时间
// 参数:
//   - spec: cron表达式，例如 "0 0 0 * * *" 表示每天0点执行
//
// 返回:
//   - time.Time: 下次执行时间
//   - error: 解析cron表达式时可能发生的错误
func GetNextTime(spec string) (time.Time, error) {
	// 创建cron解析器
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	// 解析cron表达式
	schedule, err := parser.Parse(spec)
	if err != nil {
		return time.Time{}, err
	}

	// 获取下次执行时间
	nextTime := schedule.Next(time.Now())
	return nextTime, nil
}

// GetNextTimeWithReference 根据cron表达式和参考时间获取下次执行时间
// 参数:
//   - spec: cron表达式，例如 "0 0 0 * * *" 表示每天0点执行
//   - reference: 参考时间，下次执行时间将基于此时间计算
//
// 返回:
//   - time.Time: 下次执行时间
//   - error: 解析cron表达式时可能发生的错误
func GetNextTimeWithReference(spec string, reference time.Time) (time.Time, error) {
	// 创建cron解析器
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	// 解析cron表达式
	schedule, err := parser.Parse(spec)
	if err != nil {
		return time.Time{}, err
	}

	// 获取下次执行时间
	nextTime := schedule.Next(reference)
	return nextTime, nil
}
