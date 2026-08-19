package xrcorn

import "xr-game-server/core/event"

func Init() {
	event.Sub(event.PrepareRestart, func(val any) {
		Pause()
	})
	initDay()
	initWeek()
	initMonth()
}
