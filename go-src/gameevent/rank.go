package gameevent

const (
	// AddRankValEvent 榜单数值增
	AddRankValEvent = "AddRankValEvent"
	// ReduceRankValEvent 榜单数值减
	ReduceRankValEvent = "ReduceRankValEvent"
)

type AddRankValEventData struct {
	Id     uint64
	Val    uint64
	TypeId uint32
}

func NewAddRankValEvent(roleId uint64, val uint64, typeId uint32) *AddRankValEventData {
	return &AddRankValEventData{
		Id:     roleId,
		Val:    val,
		TypeId: typeId,
	}
}

type ReduceRankEventData struct {
	Id     uint64
	Val    uint64
	TypeId uint32
}

func NewReduceRankEvent(roleId uint64, val uint64, typeId uint32) *ReduceRankEventData {
	return &ReduceRankEventData{
		Id:     roleId,
		Val:    val,
		TypeId: typeId,
	}
}
