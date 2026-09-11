package simulatordevicewhitelist

import (
	"strings"
	"sync/atomic"

	"xr-game-server/dao/simulatordevicewhitelistdao"
)

var deviceIdCache atomic.Value // map[string]struct{}

func Init() {
	reloadDeviceMemory()
}

func reloadDeviceMemory() {
	ids := simulatordevicewhitelistdao.ListAllDeviceIds()
	m := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		m[id] = struct{}{}
	}
	deviceIdCache.Store(m)
}

func getCachedDeviceSet() map[string]struct{} {
	v := deviceIdCache.Load()
	if v == nil {
		return nil
	}
	m, ok := v.(map[string]struct{})
	if !ok {
		return nil
	}
	return m
}

// IsWhitelisted 设备码是否在模拟器白名单(精确匹配)
func IsWhitelisted(deviceId string) bool {
	deviceId = strings.TrimSpace(deviceId)
	if deviceId == "" {
		return false
	}
	m := getCachedDeviceSet()
	if m == nil {
		return false
	}
	_, ok := m[deviceId]
	return ok
}
