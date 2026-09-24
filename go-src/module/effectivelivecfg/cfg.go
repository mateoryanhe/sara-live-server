package effectivelivecfg

import "xr-game-server/dao/cfgdao"

// MinSessionMinutes 返回单场直播计入有效时长的门槛(分钟).
func MinSessionMinutes() int {
	return cfgdao.EffectiveLiveMinSessionMinutes()
}

// MinSessionSeconds 返回单场直播计入有效时长的门槛(秒).
func MinSessionSeconds() float64 {
	return float64(MinSessionMinutes() * 60)
}
