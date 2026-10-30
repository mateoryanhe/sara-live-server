package liveroom

import (
	"sync"

	"xr-game-server/module/effectivelivecfg"
)

// 开播会话内摄像头开关(内存态;默认开;下播清除).是否累计有效直播时长见 ShouldAccumulateEffectiveLiveDuration.
var liveCameraEnabled sync.Map // roomId -> bool

func resetLiveCameraEnabled(roomId uint64) {
	if roomId == 0 {
		return
	}
	liveCameraEnabled.Store(roomId, true)
}

func clearLiveCameraState(roomId uint64) {
	if roomId == 0 {
		return
	}
	liveCameraEnabled.Delete(roomId)
}

func setLiveCameraEnabled(roomId uint64, enabled bool) {
	if roomId == 0 {
		return
	}
	liveCameraEnabled.Store(roomId, enabled)
}

// IsLiveCameraEnabled 未设置或下播前默认 true(摄像头开).
func IsLiveCameraEnabled(roomId uint64) bool {
	if roomId == 0 {
		return true
	}
	v, ok := liveCameraEnabled.Load(roomId)
	if !ok {
		return true
	}
	enabled, _ := v.(bool)
	return enabled
}

// ShouldAccumulateEffectiveLiveDuration 开播心跳是否累计有效直播时长(含日表/工会镜像/未结算时长).
func ShouldAccumulateEffectiveLiveDuration(roomId uint64) bool {
	if effectivelivecfg.CameraOffCountsAsEffectiveLive() {
		return true
	}
	return IsLiveCameraEnabled(roomId)
}
