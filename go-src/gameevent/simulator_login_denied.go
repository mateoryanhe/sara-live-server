package gameevent

import (
	"time"

	"xr-game-server/core/event"
	userentity "xr-game-server/entity/user"
)

const (
	// SimulatorLoginDeniedEvent 模拟器登录被拦截
	SimulatorLoginDeniedEvent event.Type = "SimulatorLoginDeniedEvent"
)

// SimulatorLoginDeniedEventData 模拟器登录拦截事件
type SimulatorLoginDeniedEventData struct {
	IP          string
	DeviceId    string
	CpuModel    string
	DeviceType  string
	DeviceModel string
	OsVersion   string
	AppVersion  string
	PackageName string
	LoginType   string
	Email       string
	AccountId   uint64
	RegisterAt  *time.Time
	RejectedAt  time.Time
}

func NewSimulatorLoginDeniedEventData(
	ip string,
	info *userentity.DeviceInfo,
	loginType, email string,
	accountId uint64,
	registerAt *time.Time,
) *SimulatorLoginDeniedEventData {
	data := &SimulatorLoginDeniedEventData{
		IP:         ip,
		LoginType:  loginType,
		Email:      email,
		AccountId:  accountId,
		RegisterAt: registerAt,
		RejectedAt: time.Now(),
	}
	if info != nil {
		data.DeviceId = info.DeviceId
		data.CpuModel = info.CpuModel
		data.DeviceType = info.DeviceType
		data.DeviceModel = info.DeviceModel
		data.OsVersion = info.OsVersion
		data.AppVersion = info.AppVersion
		data.PackageName = info.PackageName
	}
	return data
}
