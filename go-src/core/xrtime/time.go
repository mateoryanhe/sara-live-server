package xrtime

import "time"

// IsSameDay 判断两个时间是否为同一天
func IsSameDay(t1, t2 time.Time) bool {
	t1Local := t1.Local()
	t2Local := t2.Local()

	return t1Local.Year() == t2Local.Year() &&
		t1Local.Month() == t2Local.Month() &&
		t1Local.Day() == t2Local.Day()
}

// IsSameWeek 判断两个时间是否在同一周（以周一为一周的第一天）
func IsSameWeek(t1, t2 time.Time) bool {
	t1Local := t1.Local()
	t2Local := t2.Local()

	weekStart1 := getWeekStart(t1Local)
	weekStart2 := getWeekStart(t2Local)

	return weekStart1.Equal(weekStart2)
}

// 获取时间所在周的周一
func getWeekStart(t time.Time) time.Time {
	weekday := t.Weekday()
	offset := int(weekday - time.Monday)
	if offset < 0 {
		offset += 7
	}
	return time.Date(t.Year(), t.Month(), t.Day()-offset, 0, 0, 0, 0, t.Location())
}

// IsSameMonth 判断两个时间是否属于同一个月
func IsSameMonth(t1, t2 time.Time) bool {
	// 转换到同一时区（这里使用本地时区，也可以根据需要使用UTC等）
	t1Local := t1.Local()
	t2Local := t2.Local()

	// 比较年份和月份是否相同
	return t1Local.Year() == t2Local.Year() &&
		t1Local.Month() == t2Local.Month()
}
