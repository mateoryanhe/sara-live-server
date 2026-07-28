package agora

import (
	"sync"
	"time"
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
	cloudPlayerUpdateSeq.Store(playerId, int(time.Now().Unix()))
}

func nextCloudPlayerSequence(playerId string) int {
	if playerId == "" {
		return 0
	}
	val, _ := cloudPlayerUpdateSeq.LoadOrStore(playerId, int(time.Now().Unix())-1)
	seq := val.(int) + 1
	cloudPlayerUpdateSeq.Store(playerId, seq)
	return seq
}
