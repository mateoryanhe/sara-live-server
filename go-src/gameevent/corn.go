package gameevent

import "xr-game-server/core/event"

const (
	DayEvent   event.Type = "DayEvent"
	WeekEvent  event.Type = "WeekEvent"
	MonthEvent event.Type = "MonthEvent"
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
