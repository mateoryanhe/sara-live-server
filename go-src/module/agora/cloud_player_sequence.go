package agora

import (
	"sync"
)

var cloudPlayerUpdateSeq sync.Map // playerId -> int

func registerCloudPlayerSequence(playerId string) {
	if playerId == "" {
		return
	}
	cloudPlayerUpdateSeq.Store(playerId, -1)
}

func unregisterCloudPlayerSequence(playerId string) {
	if playerId == "" {
		return
	}
	cloudPlayerUpdateSeq.Delete(playerId)
}

func restoreCloudPlayerSequence(playerId string) {
	if playerId == "" {
		return
	}
	if _, loaded := cloudPlayerUpdateSeq.Load(playerId); loaded {
		return
	}
	// 服务重启后从 0 开始；若 Agora 侧已有更高 sequence，续期失败会走重建逻辑。
	cloudPlayerUpdateSeq.Store(playerId, -1)
}

func nextCloudPlayerSequence(playerId string) int {
	if playerId == "" {
		return 0
	}
	val, _ := cloudPlayerUpdateSeq.LoadOrStore(playerId, -1)
	seq := val.(int) + 1
	cloudPlayerUpdateSeq.Store(playerId, seq)
	return seq
}
