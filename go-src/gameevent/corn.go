package gameevent

import "xr-game-server/core/event"

const (
	DayEvent          event.Type = "DayEvent"
	WeekEvent         event.Type = "WeekEvent"
	MonthEvent        event.Type = "MonthEvent"
	TimezoneDayEvent  event.Type = "TimezoneDayEvent"
	TimezoneWeekEvent event.Type = "TimezoneWeekEvent"
)

type DayEventData struct {
}

func NewDayEventData() *DayEventData {
	return &DayEventData{}
}

type WeekEventData struct {
}

func NewWeekEventData() *WeekEventData {
	return &WeekEventData{}
}

type MonthEventData struct {
}

func NewMonthEventData() *MonthEventData {
	return &MonthEventData{}
}

type TimezoneDayEventData struct {
	Timezone int `json:"timezone"`
}

func NewTimezoneDayEventData(timezone int) *TimezoneDayEventData {
	return &TimezoneDayEventData{Timezone: timezone}
}

type TimezoneWeekEventData struct {
	Timezone int `json:"timezone"`
}

func NewTimezoneWeekEventData(timezone int) *TimezoneWeekEventData {
	return &TimezoneWeekEventData{Timezone: timezone}
}

func InitTimezoneEvents() {
	// 预留扩展点，时区事件初始化（如需要动态订阅）
}
