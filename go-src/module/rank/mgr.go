package rank

import (
	"xr-game-server/core/actor"
)

// 榜单管理器
var mgrMap = make(map[uint32]*Mgr, 100)

// 榜单数据缓冲
var rankActor = actor.NewActor(10000)

type Mgr struct {
	ZoneRank   map[uint64]*CommonRank
	ServerRank map[uint64]*CommonRank
}

func NewMgr() *Mgr {
	return &Mgr{
		ZoneRank:   make(map[uint64]*CommonRank),
		ServerRank: make(map[uint64]*CommonRank),
	}
}

func (receiver *Mgr) GetRankByServerId(serverId uint64) *CommonRank {
	//serverData := serverdao.GetServer(serverId)
	//if serverData.ZoneId > common.Zero {
	//	return receiver.ZoneRank[serverData.ZoneId]
	//} else {
	//	return receiver.ServerRank[serverId]
	//}
	return receiver.ServerRank[serverId]
}

func GetRankBy(typeId uint32, serverId uint64) *CommonRank {
	//mgr, mgrOk := mgrMap[typeId]
	//if !mgrOk {
	//	return nil
	//}
	//ret, serverOk := mgr.ServerRank[serverId]
	//if serverOk {
	//	return ret
	//}
	//server := serverdao.GetServer(serverId)
	//zoneRank, zoneOk := mgr.ZoneRank[server.ZoneId]
	//if zoneOk {
	//	return zoneRank
	//}
	return nil
}

func initMgr() {
	rankActor.Start()
}

func onAddEvent(event any) {
	rankActor.Send(event, func(msg any) {
		//data := msg.(*gameevent.AddRankValEventData)
		//mgr := mgrMap[data.TypeId]
		//roleData := roledao.GetRoleByRole(data.Id)
		//commonRank := mgr.GetRankByServerId(roleData.ServerId)
		//commonRank.Add(data)
	})

}
func onReduceEvent(event any) {
	rankActor.Send(event, func(msg any) {
		//data := msg.(*gameevent.ReduceRankEventData)
		//mgr := mgrMap[data.TypeId]
		//roleData := roledao.GetRoleByRole(data.Id)
		//commonRank := mgr.GetRankByServerId(roleData.ServerId)
		//commonRank.Reduce(data)
	})
}
