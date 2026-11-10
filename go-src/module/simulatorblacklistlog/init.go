package simulatorblacklistlog

import (
	"xr-game-server/core/event"
	"xr-game-server/gameevent"
	userentity "xr-game-server/entity/user"
)

func Init() {
	event.Sub(gameevent.SimulatorLoginDeniedEvent, onSimulatorLoginDenied)
}

func onSimulatorLoginDenied(data any) {
	ev, ok := data.(*gameevent.SimulatorLoginDeniedEventData)
	if !ok || ev == nil {
		return
	}
	_ = userentity.NewSimulatorBlacklistLog(
		ev.IP,
		ev.DeviceId,
		ev.CpuModel,
		ev.DeviceType,
		ev.DeviceModel,
		ev.OsVersion,
		ev.AppVersion,
		ev.PackageName,
		ev.LoginType,
		ev.Email,
		ev.AccountId,
		ev.RegisterAt,
	)
}
