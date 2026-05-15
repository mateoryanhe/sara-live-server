package rank

import (
	"xr-game-server/core/event"
	"xr-game-server/gameevent"
)

// 天,周,月结算
func initSettlement() {
	event.Sub(gameevent.DayEvent, initDay)
	event.Sub(gameevent.WeekEvent, initWeek)
	event.Sub(gameevent.MonthEvent, initMonth)
}

func initDay(data any) {

}

func initWeek(data any) {

}

func initMonth(data any) {

}
