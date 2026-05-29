package liveroom

import "xr-game-server/core/actor"

const (
	Max  = 5
	Size = 100
)

var actorMap = make(map[uint64]*actor.Actor)

func initRoomActor() {
	for i := 0; i < Max; i++ {
		val := actor.NewActor(Size)
		val.Start()
		actorMap[uint64(i)] = val
	}
}
